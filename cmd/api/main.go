package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/davidcm146/assets-management-be.git/internal/config"
	"github.com/davidcm146/assets-management-be.git/internal/database"
	"github.com/davidcm146/assets-management-be.git/internal/handler"
	"github.com/davidcm146/assets-management-be.git/internal/infrastructure/cloudinary"
	"github.com/davidcm146/assets-management-be.git/internal/infrastructure/gmail"
	"github.com/davidcm146/assets-management-be.git/internal/infrastructure/postmark"
	"github.com/davidcm146/assets-management-be.git/internal/mailer"
	"github.com/davidcm146/assets-management-be.git/internal/repository"
	"github.com/davidcm146/assets-management-be.git/internal/router"
	"github.com/davidcm146/assets-management-be.git/internal/scheduler"
	"github.com/davidcm146/assets-management-be.git/internal/scheduler/jobs"
	"github.com/davidcm146/assets-management-be.git/internal/server"
	"github.com/davidcm146/assets-management-be.git/internal/service"
	"github.com/davidcm146/assets-management-be.git/internal/validator"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := validator.RegisterValidators(); err != nil {
		panic(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := database.NewDB(ctx, cfg.DB)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	cld, _ := cloudinary.NewCloudinary(&cfg.Cloudinary)
	postmarkProvider := postmark.NewProvider(&cfg.Postmark)
	gmailProvider := gmail.NewProvider(&cfg.Gmail)
	uploader := cloudinary.NewCloudinaryUploader(cld)
	mailProvider := service.NewMultiMailProvider(postmarkProvider, gmailProvider)
	renderer, err := email.NewRenderer()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	loanSlipRepo := repository.NewLoanSlipRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)

	authService := service.NewAuthService(userRepo)
	loanSlipService := service.NewLoanSlipService(loanSlipRepo, uploader)
	notificationService := service.NewNotificationService(renderer, notificationRepo, mailProvider)
	dashboardService := service.NewDashboardService(dashboardRepo)

	authHandler := handler.NewAuthHandler(authService, userRepo)
	loanSlipHandler := handler.NewLoanSlipHandler(loanSlipService, uploader)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	engine := gin.New()
	r := router.NewRouter(router.RouterParams{
		Engine: engine,
		Handlers: &router.Handlers{
			AuthHandler:         authHandler,
			LoanSlipHandler:     loanSlipHandler,
			NotificationHandler: notificationHandler,
			DashboardHandler:    dashboardHandler,
		},
	})
	srv := server.NewServer(r, cfg.Server.Port)

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	cronScheduler := scheduler.NewScheduler()
	overdueJob := jobs.NewOverdueJob(loanSlipService, notificationService)

	if err := scheduler.RegisterJobs(cronScheduler, []jobs.Job{
		overdueJob,
	}); err != nil {
		panic(err)
	}
	cronScheduler.Start()

	<-ctx.Done()
	cronScheduler.Stop()
	srv.Shutdown(context.Background())
}

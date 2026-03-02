package router

import (
	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/handler"
	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/davidcm146/assets-management-be.git/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	AuthHandler         handler.AuthHandler
	LoanSlipHandler     handler.LoanSlipHandler
	NotificationHandler handler.NotificationHandler
	DashboardHandler    handler.DashboardHandler
}

type RouterParams struct {
	Engine   *gin.Engine
	Handlers *Handlers
}

func NewRouter(params RouterParams) *gin.Engine {
	r := params.Engine

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(error_middleware.ErrorHandler())

	// Routes
	api := r.Group("/api")
	auth_api := api.Group("/auth")
	{
		auth_api.POST("/register", params.Handlers.AuthHandler.RegisterHandler)
		auth_api.POST("/login", params.Handlers.AuthHandler.LoginHandler)
	}

	protected_api := api.Group("/")
	protected_api.Use(middleware.AuthMiddleware())
	{
		protected_api.GET("/me", params.Handlers.AuthHandler.GetMe)
		protected_api.GET("/loan-slips", params.Handlers.LoanSlipHandler.LoanSlipsListHandler)
		protected_api.GET("/loan-slips/:id", params.Handlers.LoanSlipHandler.LoanSlipDetailHandler)
		protected_api.POST("/loan-slips", middleware.PermissionMiddleware([]model.Role{model.Admin, model.IT}), params.Handlers.LoanSlipHandler.CreateLoanSlipHandler)
		protected_api.PUT("/loan-slips/:id", middleware.PermissionMiddleware([]model.Role{model.Admin, model.IT}), params.Handlers.LoanSlipHandler.UpdateLoanSlipHandler)
		protected_api.DELETE("/loan-slips/:id", middleware.PermissionMiddleware([]model.Role{model.Admin, model.IT}), params.Handlers.LoanSlipHandler.DeleteLoanSlipHandler)
		protected_api.GET("/notifications", middleware.PermissionMiddleware([]model.Role{model.Admin, model.IT}), params.Handlers.NotificationHandler.ListHandler)
		protected_api.PUT("/notifications/:id", middleware.PermissionMiddleware([]model.Role{model.Admin, model.IT}), params.Handlers.NotificationHandler.MarkAsReadHandler)
		protected_api.GET("/notifications/unread/count", middleware.PermissionMiddleware([]model.Role{model.Admin, model.IT}), params.Handlers.NotificationHandler.CountUnreadHandler)
	}
	dashboard := protected_api.Group("/dashboard")
	{
		dashboard.GET("/loan-metrics", middleware.PermissionMiddleware([]model.Role{model.Admin, model.IT}), params.Handlers.DashboardHandler.GetLoanMetricsHandler)
	}

	return r
}

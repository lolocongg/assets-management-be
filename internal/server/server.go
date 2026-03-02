package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(r *gin.Engine, port string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: r,
		},
	}
}

func (s *Server) Run() error {
	fmt.Println("Server listening on", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}

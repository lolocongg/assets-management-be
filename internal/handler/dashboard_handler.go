package handler

import (
	"net/http"

	"github.com/davidcm146/assets-management-be.git/internal/dto"
	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/service"
	"github.com/gin-gonic/gin"
)

type DashboardHandler interface {
	GetLoanMetricsHandler(c *gin.Context)
}

type dashboardHandler struct {
	dashboardService service.DashboardService
}

func NewDashboardHandler(dashboardService service.DashboardService) DashboardHandler {
	return &dashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *dashboardHandler) GetLoanMetricsHandler(c *gin.Context) {
	var req dto.DashboardFilterRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(error_middleware.NewBadRequest("Query không hợp lệ"))
		return
	}

	result, err := h.dashboardService.GetLoanMetricsService(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

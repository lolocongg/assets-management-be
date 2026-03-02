package handler

import (
	"net/http"
	"strconv"

	"github.com/davidcm146/assets-management-be.git/internal/dto"
	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/service"
	"github.com/gin-gonic/gin"
)

type NotificationHandler interface {
	ListHandler(c *gin.Context)
	MarkAsReadHandler(c *gin.Context)
	CountUnreadHandler(c *gin.Context)
}

type notificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) NotificationHandler {
	return &notificationHandler{notificationService: notificationService}
}

func (h *notificationHandler) ListHandler(c *gin.Context) {
	user := c.MustGet("user").(*dto.AuthUser)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, total, err := h.notificationService.List(c.Request.Context(), user.ID, page, limit, nil)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": total,
	})
}

func (h *notificationHandler) MarkAsReadHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(error_middleware.NewBadRequest("ID không hợp lệ"))
		return
	}

	err = h.notificationService.MarkAsRead(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã đánh dấu đã đọc"})
}

func (h *notificationHandler) CountUnreadHandler(c *gin.Context) {
	user := c.MustGet("user").(*dto.AuthUser)

	count, err := h.notificationService.CountUnread(c.Request.Context(), user.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"unread_count": count,
	})
}

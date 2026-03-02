package handler

import (
	"net/http"
	"strconv"

	"github.com/davidcm146/assets-management-be.git/internal/dto"
	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/policy"
	"github.com/davidcm146/assets-management-be.git/internal/service"
	"github.com/davidcm146/assets-management-be.git/internal/validator"
	"github.com/gin-gonic/gin"
)

type LoanSlipHandler interface {
	LoanSlipsListHandler(c *gin.Context)
	CreateLoanSlipHandler(c *gin.Context)
	UpdateLoanSlipHandler(c *gin.Context)
	LoanSlipDetailHandler(c *gin.Context)
	DeleteLoanSlipHandler(c *gin.Context)
}

type loanSlipHandler struct {
	loanSlipService service.LoanSlipService
	uploader        service.Uploader
	policy          policy.LoanSlipPolicy
}

func NewLoanSlipHandler(loanSlipService service.LoanSlipService, uploader service.Uploader) LoanSlipHandler {
	return &loanSlipHandler{
		loanSlipService: loanSlipService,
		uploader:        uploader,
	}
}

func (h *loanSlipHandler) LoanSlipsListHandler(c *gin.Context) {
	var query dto.LoanSlipQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(error_middleware.NewBadRequest("Yêu cầu không hợp lệ"))
		return
	}
	result, err := h.loanSlipService.LoanSlipsListService(c.Request.Context(), &query)

	if err != nil {
		if _, ok := err.(*error_middleware.AppError); ok {
			c.Error(err)
			return
		}

		c.Error(error_middleware.NewInternal("Lỗi hệ thống"))
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *loanSlipHandler) CreateLoanSlipHandler(c *gin.Context) {
	var req dto.CreateLoanSlipRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(
			error_middleware.NewUnprocessableEntity("Dữ liệu không hợp lệ").
				WithDetails(validator.HandleValidationError(err, &req)),
		)
		return
	}

	userRaw, exists := c.Get("user")
	if !exists {
		c.Error(error_middleware.NewUnauthorized("Unauthorized"))
		return
	}
	user := userRaw.(*dto.AuthUser)

	loan, err := h.loanSlipService.CreateLoanSlipService(c.Request.Context(), user.ID, &req)

	if err != nil {
		if _, ok := err.(*error_middleware.AppError); ok {
			c.Error(err)
			return
		}

		c.Error(error_middleware.NewInternal("Lỗi hệ thống"))
		return
	}

	c.JSON(http.StatusCreated, loan)
}

func extractUpdateFields(updateDTO *dto.UpdateLoanSlipRequest) []string {
	fields := []string{}

	if updateDTO.Name != nil {
		fields = append(fields, "name")
	}
	if updateDTO.BorrowerName != nil {
		fields = append(fields, "borrower_name")
	}
	if updateDTO.Department != nil {
		fields = append(fields, "department")
	}
	if updateDTO.Position != nil {
		fields = append(fields, "position")
	}
	if updateDTO.Description != nil {
		fields = append(fields, "description")
	}
	if updateDTO.Status != nil {
		fields = append(fields, "status")
	}
	if updateDTO.SerialNumber != nil {
		fields = append(fields, "serial_number")
	}
	if updateDTO.BorrowedDate != nil {
		fields = append(fields, "borrowed_date")
	}
	if updateDTO.ReturnedDate != nil {
		fields = append(fields, "returned_date")
	}
	return fields
}

func (h *loanSlipHandler) UpdateLoanSlipHandler(c *gin.Context) {
	var req dto.UpdateLoanSlipRequest

	if err := c.ShouldBind(&req); err != nil {
		c.Error(
			error_middleware.NewUnprocessableEntity("Dữ liệu không hợp lệ").
				WithDetails(validator.HandleValidationError(err, &req)),
		)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(error_middleware.NewBadRequest("ID không hợp lệ"))
		return
	}

	user := c.MustGet("user").(*dto.AuthUser)
	fields := extractUpdateFields(&req)

	if forbidden := h.policy.ForbiddenFields(user.Role, fields); len(forbidden) > 0 {
		c.Error(
			error_middleware.
				NewForbidden("Bạn không có quyền cập nhật một số trường").
				WithDetails(map[string]any{
					"fields": forbidden,
				}),
		)
		return
	}

	loanSlip, err := h.loanSlipService.UpdateLoanSlipService(c.Request.Context(), id, &req)
	if err != nil {
		if _, ok := err.(*error_middleware.AppError); ok {
			c.Error(err)
			return
		}

		c.Error(error_middleware.NewInternal("Lỗi hệ thống"))
		return
	}

	c.JSON(http.StatusOK, loanSlip)
}

func (h *loanSlipHandler) DeleteLoanSlipHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(error_middleware.NewBadRequest("ID không hợp lệ"))
		return
	}

	if err = h.loanSlipService.Delete(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *loanSlipHandler) LoanSlipDetailHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(error_middleware.NewBadRequest("ID không hợp lệ"))
		return
	}

	result, err := h.loanSlipService.LoanSlipDetailService(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

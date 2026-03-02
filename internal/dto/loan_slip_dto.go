package dto

import (
	"mime/multipart"
	"time"

	"github.com/davidcm146/assets-management-be.git/internal/model"
)

type LoanSlipResponse struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	BorrowerName string     `json:"borrower_name"`
	Department   string     `json:"department"`
	Position     string     `json:"position"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	SerialNumber string     `json:"serial_number"`
	Images       []string   `json:"images"`
	BorrowedDate *time.Time `json:"borrowed_date"`
	ReturnedDate *time.Time `json:"returned_date"`
	CreatedAt    *time.Time `json:"created_at"`
}

type LoanSlipQuery struct {
	Search string `form:"search"`

	Department string `form:"department"`
	Status     string `form:"status"`

	BorrowedFrom *time.Time `form:"borrowed_from" time_format:"02-01-2006"`
	BorrowedTo   *time.Time `form:"borrowed_to" time_format:"02-01-2006"`

	ReturnedFrom *time.Time `form:"returned_from" time_format:"02-01-2006"`
	ReturnedTo   *time.Time `form:"returned_to" time_format:"02-01-2006"`

	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=10"`

	Sort  string `form:"sort"`
	Order string `form:"order"`
}

type PagedResult[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
}

func NewPagedResult[T any](items []T, total int) *PagedResult[T] {
	return &PagedResult[T]{
		Items: items,
		Total: total,
	}
}

type CreateLoanSlipRequest struct {
	Name         string                  `form:"name" binding:"required,gte=3,lte=100" label:"Tên tài sản"`
	BorrowerName string                  `form:"borrower_name" binding:"required,gte=3,lte=50" label:"Tên nhà thầu"`
	Department   string                  `form:"department" label:"Phòng ban"`
	Position     string                  `form:"position" label:"Chức vụ"`
	Description  string                  `form:"description" binding:"lte=500" label:"Mô tả"`
	SerialNumber string                  `form:"serial_number" label:"Số sê ri"`
	BorrowedDate *time.Time              `form:"borrowed_date" binding:"required" time_format:"02-01-2006" label:"Ngày mượn"`
	ReturnedDate *time.Time              `form:"returned_date" binding:"required,gtefield=BorrowedDate" time_format:"02-01-2006" label:"Ngày trả"`
	Images       []*multipart.FileHeader `form:"images" binding:"omitempty,max=5,images" label:"Hình ảnh"`
}

type UpdateLoanSlipRequest struct {
	Name         *string                 `form:"name" label:"Tên tài sản"`
	BorrowerName *string                 `form:"borrower_name" label:"Tên nhà thầu"`
	Department   *string                 `form:"department" label:"Phòng ban"`
	Position     *string                 `form:"position" label:"Chức vụ"`
	Description  *string                 `form:"description" label:"Mô tả"`
	SerialNumber *string                 `form:"serial_number" label:"Số sê ri"`
	Status       *model.Status           `form:"status" label:"Trạng thái"`
	BorrowedDate *time.Time              `form:"borrowed_date" time_format:"02-01-2006" label:"Ngày mượn"`
	ReturnedDate *time.Time              `form:"returned_date" binding:"gtefield=BorrowedDate" time_format:"02-01-2006" label:"Ngày trả"`
	Images       []*multipart.FileHeader `form:"images" binding:"omitempty,max=5,images" label:"Hình ảnh"`
}

package service

import (
	"context"
	"mime/multipart"

	"github.com/davidcm146/assets-management-be.git/internal/dto"
	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/davidcm146/assets-management-be.git/internal/policy"
	"github.com/davidcm146/assets-management-be.git/internal/repository"
	"golang.org/x/sync/errgroup"
)

type LoanSlipService interface {
	LoanSlipsListService(ctx context.Context, query *dto.LoanSlipQuery) (*dto.PagedResult[*model.LoanSlip], error)
	uploadImages(ctx context.Context, files []*multipart.FileHeader) ([]string, error)
	CreateLoanSlipService(ctx context.Context, userID int, req *dto.CreateLoanSlipRequest) (*model.LoanSlip, error)
	UpdateLoanSlipService(ctx context.Context, id int, updateDTO *dto.UpdateLoanSlipRequest) (*model.LoanSlip, error)
	LoanSlipDetailService(ctx context.Context, id int) (*dto.LoanSlipResponse, error)
	MarkAsOverdue(ctx context.Context, id int) (*model.LoanSlip, error)
	MarkOverdueNotified(ctx context.Context, id int) (*model.LoanSlip, error)
	GetOverdue(ctx context.Context) ([]*model.LoanSlip, error)
	Delete(ctx context.Context, id int) error
}

type loanSlipService struct {
	loanSlipRepo repository.LoanSlipRepository
	uploader     Uploader
	policy       policy.LoanSlipPolicy
}

func NewLoanSlipService(loanSlipRepo repository.LoanSlipRepository, uploader Uploader) LoanSlipService {
	return &loanSlipService{loanSlipRepo: loanSlipRepo, uploader: uploader}
}

func mapLoanSlipToResponse(m *model.LoanSlip) *dto.LoanSlipResponse {
	return &dto.LoanSlipResponse{
		ID:           m.ID,
		Name:         m.Name,
		BorrowerName: m.BorrowerName,
		Department:   m.Department,
		Position:     m.Position,
		Description:  m.Description,
		Status:       m.Status.String(),
		SerialNumber: m.SerialNumber,
		Images:       m.Images,
		BorrowedDate: m.BorrowedDate,
		ReturnedDate: m.ReturnedDate,
		CreatedAt:    m.CreatedAt,
	}
}

func (s *loanSlipService) LoanSlipsListService(ctx context.Context, query *dto.LoanSlipQuery) (*dto.PagedResult[*model.LoanSlip], error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}
	if query.Sort == "" {
		query.Sort = "created_at"
	}
	if query.Order != "ASC" {
		query.Order = "DESC"
	}

	items, err := s.loanSlipRepo.List(ctx, query)
	if err != nil {
		return nil, error_middleware.NewInternal("Không thể lấy danh sách phiếu mượn")
	}

	if items == nil {
		items = []*model.LoanSlip{}
	}
	total, err := s.loanSlipRepo.Count(ctx, query)
	if err != nil {
		return nil, error_middleware.NewInternal("Không thể đếm số lượng phiếu mượn")
	}

	return dto.NewPagedResult(items, total), nil
}

func (s *loanSlipService) uploadImages(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {

	g, ctx := errgroup.WithContext(ctx)

	imageURLs := make([]string, len(files))

	for i, file := range files {
		g.Go(func() error {
			url, err := s.uploader.Upload(ctx, file)
			if err != nil {
				return err
			}

			imageURLs[i] = url
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return imageURLs, nil
}

func (s *loanSlipService) CreateLoanSlipService(ctx context.Context, userID int, req *dto.CreateLoanSlipRequest) (*model.LoanSlip, error) {
	urls, err := s.uploadImages(ctx, req.Images)
	if err != nil {
		return nil, error_middleware.NewInternal("Lỗi upload ảnh")
	}

	loan := &model.LoanSlip{
		Name:         req.Name,
		BorrowerName: req.BorrowerName,
		Department:   req.Department,
		Position:     req.Position,
		Description:  req.Description,
		SerialNumber: req.SerialNumber,
		Images:       urls,
		CreatedBy:    userID,
		BorrowedDate: req.BorrowedDate,
		ReturnedDate: req.ReturnedDate,
		Status:       model.Borrowing,
	}

	if err := s.loanSlipRepo.Create(ctx, loan); err != nil {
		return nil, error_middleware.NewInternal("Tạo phiếu mượn thất bại")
	}

	return loan, nil
}

func applyLoanSlipUpdate(loanSlip *model.LoanSlip, updateDTO *dto.UpdateLoanSlipRequest) {
	if updateDTO.Name != nil {
		loanSlip.Name = *updateDTO.Name
	}
	if updateDTO.BorrowerName != nil {
		loanSlip.BorrowerName = *updateDTO.BorrowerName
	}
	if updateDTO.Department != nil {
		loanSlip.Department = *updateDTO.Department
	}
	if updateDTO.Position != nil {
		loanSlip.Position = *updateDTO.Position
	}
	if updateDTO.Description != nil {
		loanSlip.Description = *updateDTO.Description
	}
	if updateDTO.Status != nil {
		loanSlip.Status = model.Status(*updateDTO.Status)
	}
	if updateDTO.SerialNumber != nil {
		loanSlip.SerialNumber = *updateDTO.SerialNumber
	}
	if updateDTO.BorrowedDate != nil {
		loanSlip.BorrowedDate = updateDTO.BorrowedDate
	}
	if updateDTO.ReturnedDate != nil {
		loanSlip.ReturnedDate = updateDTO.ReturnedDate
	}
}

func (s *loanSlipService) UpdateLoanSlipService(ctx context.Context, id int, updateDTO *dto.UpdateLoanSlipRequest) (*model.LoanSlip, error) {
	loanSlip, err := s.loanSlipRepo.FindByID(ctx, id)
	if err != nil {
		return nil, error_middleware.NewNotFound("Không tìm thấy phiếu mượn")
	}

	if updateDTO.Status != nil {
		newStatus := model.Status(*updateDTO.Status)

		if !loanSlip.Status.CanTransition(newStatus) {
			return nil, error_middleware.NewUnprocessableEntity("Không thể chuyển trạng thái từ đã trả sang đang mượn")
		}
	}

	if len(updateDTO.Images) > 0 {
		urls, err := s.uploadImages(ctx, updateDTO.Images)
		if err != nil {
			return nil, error_middleware.NewInternal("Lỗi upload ảnh")
		}
		loanSlip.Images = append(loanSlip.Images, urls...)
	}

	applyLoanSlipUpdate(loanSlip, updateDTO)

	if err := s.loanSlipRepo.Update(ctx, loanSlip); err != nil {
		return nil, error_middleware.NewInternal("Cập nhật thất bại")
	}

	return loanSlip, nil
}

func (s *loanSlipService) Delete(ctx context.Context, id int) error {
	loanSlip, err := s.loanSlipRepo.FindByID(ctx, id)
	if err != nil || loanSlip == nil {
		return error_middleware.NewNotFound("Không tìm thấy phiếu mượn")
	}

	if loanSlip.Status != model.Borrowing {
		return error_middleware.NewUnprocessableEntity("Không thể xóa phiếu mượn ở trạng thái hiện tại")
	}
	if err := s.loanSlipRepo.Delete(ctx, id); err != nil {
		return error_middleware.NewInternal("Xóa phiếu mượn thất bại")
	}

	return nil
}

func (s *loanSlipService) LoanSlipDetailService(ctx context.Context, id int) (*dto.LoanSlipResponse, error) {
	loanSlip, err := s.loanSlipRepo.FindByID(ctx, id)
	if err != nil {
		return nil, error_middleware.NewNotFound("Không tìm thấy phiếu mượn")
	}
	return mapLoanSlipToResponse(loanSlip), nil
}

func (s *loanSlipService) MarkAsOverdue(ctx context.Context, id int) (*model.LoanSlip, error) {

	loanSlip, err := s.loanSlipRepo.FindByID(ctx, id)
	if err != nil {
		return nil, error_middleware.NewNotFound("Không tìm thấy phiếu mượn")
	}
	if loanSlip.Status == model.Returned {
		return nil, error_middleware.NewUnprocessableEntity("Phiếu mượn đã trả")
	}

	loanSlip.Status = model.Overdue

	if err := s.loanSlipRepo.Update(ctx, loanSlip); err != nil {
		return nil, error_middleware.NewInternal("Cập nhật trạng thái quá hạn thất bại")
	}

	return loanSlip, nil
}

func (s *loanSlipService) MarkOverdueNotified(ctx context.Context, id int) (*model.LoanSlip, error) {
	loanSlip, err := s.loanSlipRepo.FindByID(ctx, id)
	if err != nil {
		return nil, error_middleware.NewNotFound("Không tìm thấy phiếu mượn")
	}

	if err := s.loanSlipRepo.MarkOverdueNotified(ctx, id); err != nil {
		return nil, error_middleware.NewInternal("Cập nhật overdue_notified thất bại")
	}

	loanSlip.OverdueNotified = true

	return loanSlip, nil
}

func (s *loanSlipService) GetOverdue(ctx context.Context) ([]*model.LoanSlip, error) {
	loanSlips, err := s.loanSlipRepo.FindOverdue(ctx)
	if err != nil {
		return nil, error_middleware.NewInternal("Lấy danh sách phiếu mượn quá hạn thất bại")
	}

	if loanSlips == nil {
		loanSlips = []*model.LoanSlip{}
	}

	return loanSlips, nil
}

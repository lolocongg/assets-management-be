package service

import (
	"context"
	"fmt"
	"time"

	"github.com/davidcm146/assets-management-be.git/internal/dto"
	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/davidcm146/assets-management-be.git/internal/repository"
)

type DashboardService interface {
	GetLoanMetricsService(ctx context.Context, req dto.DashboardFilterRequest) (*dto.LoanMetricsResponse, error)
}

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
}

func NewDashboardService(dashboardRepo repository.DashboardRepository) DashboardService {
	return &dashboardService{
		dashboardRepo: dashboardRepo,
	}
}

func (s *dashboardService) resolveTimeFilter(req dto.DashboardFilterRequest) (model.TimeFilter, error) {
	var filter model.TimeFilter
	layout := "02-01-2006"

	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.Local
	}

	if req.From != "" && req.To != "" {
		from, err := time.ParseInLocation(layout, req.From, loc)
		if err != nil {
			return filter, fmt.Errorf("Ngày bắt đầu không hợp lệ")
		}

		to, err := time.ParseInLocation(layout, req.To, loc)
		if err != nil {
			return filter, fmt.Errorf("Ngày đến không hợp lệ")
		}

		filter.From = &from
		filter.To = &to
		return filter, nil
	}

	now := time.Now().In(loc)

	switch req.Period {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		filter.From = &start
		filter.To = &now

	case "month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		filter.From = &start
		filter.To = &now

	case "year":
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc)
		filter.From = &start
		filter.To = &now
	}

	return filter, nil
}

func mapLoanMetricsToResponse(m *model.LoanMetrics) *dto.LoanMetricsResponse {
	if m == nil {
		return &dto.LoanMetricsResponse{}
	}

	return &dto.LoanMetricsResponse{
		Total:     m.Total,
		Borrowing: m.Borrowing,
		Returned:  m.Returned,
		Overdue:   m.Overdue,
	}
}

func (s *dashboardService) GetLoanMetricsService(ctx context.Context, req dto.DashboardFilterRequest) (*dto.LoanMetricsResponse, error) {
	filter, err := s.resolveTimeFilter(req)
	if err != nil {
		return nil, error_middleware.NewBadRequest(err.Error())
	}

	metrics, err := s.dashboardRepo.GetLoanMetrics(ctx, filter)
	if err != nil {
		return nil, error_middleware.NewInternal("Không thể lấy thống kê phiếu mượn")
	}

	return mapLoanMetricsToResponse(metrics), nil
}

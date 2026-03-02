package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository interface {
	GetLoanMetrics(ctx context.Context, filter model.TimeFilter) (*model.LoanMetrics, error)
}

type dashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetLoanMetrics(ctx context.Context, filter model.TimeFilter) (*model.LoanMetrics, error) {
	var result model.LoanMetrics

	whereBorrow := []string{"status = $STATUS_BORROWING"}
	whereReturned := []string{"status = $STATUS_RETURNED"}
	whereOverdue := []string{"status = $STATUS_OVERDUE"}

	args := []any{}
	argPos := 1

	borrowingPos := argPos
	args = append(args, model.Borrowing)
	argPos++

	returnedPos := argPos
	args = append(args, model.Returned)
	argPos++

	overduePos := argPos
	args = append(args, model.Overdue)
	argPos++
	if filter.From != nil {
		whereBorrow = append(whereBorrow,
			fmt.Sprintf("created_at >= $%d", argPos))
		whereReturned = append(whereReturned,
			fmt.Sprintf("returned_at >= $%d", argPos))
		whereOverdue = append(whereOverdue,
			fmt.Sprintf("overdue_at >= $%d", argPos))

		args = append(args, *filter.From)
		argPos++
	}

	if filter.To != nil {
		whereBorrow = append(whereBorrow,
			fmt.Sprintf("created_at <= $%d", argPos))
		whereReturned = append(whereReturned,
			fmt.Sprintf("returned_at <= $%d", argPos))
		whereOverdue = append(whereOverdue,
			fmt.Sprintf("overdue_at <= $%d", argPos))

		args = append(args, *filter.To)
		argPos++
	}

	query := fmt.Sprintf(`
		SELECT
			(SELECT COUNT(*) FROM loan_slips WHERE %s) AS borrowing,
			(SELECT COUNT(*) FROM loan_slips WHERE %s) AS returned,
			(SELECT COUNT(*) FROM loan_slips WHERE %s) AS overdue
	`,
		strings.ReplaceAll(strings.Join(whereBorrow, " AND "),
			"$STATUS_BORROWING", fmt.Sprintf("$%d", borrowingPos)),
		strings.ReplaceAll(strings.Join(whereReturned, " AND "),
			"$STATUS_RETURNED", fmt.Sprintf("$%d", returnedPos)),
		strings.ReplaceAll(strings.Join(whereOverdue, " AND "),
			"$STATUS_OVERDUE", fmt.Sprintf("$%d", overduePos)),
	)

	row := r.db.QueryRow(ctx, query, args...)
	err := row.Scan(
		&result.Borrowing,
		&result.Returned,
		&result.Overdue,
	)
	if err != nil {
		return nil, err
	}
	result.Total = result.Borrowing + result.Returned + result.Overdue

	return &result, nil
}

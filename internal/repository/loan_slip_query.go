package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/davidcm146/assets-management-be.git/internal/dto"
	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/davidcm146/assets-management-be.git/internal/utils"
)

func buildLoanSlipQuery(query *dto.LoanSlipQuery) (string, []any) {
	conds := []string{"1=1"}
	args := []any{}
	argIdx := 1

	if query.Search != "" {
		conds = append(conds, fmt.Sprintf(
			"(name ILIKE $%d OR department ILIKE $%d OR borrower_name ILIKE $%d)",
			argIdx, argIdx+1, argIdx+2,
		))
		args = append(args,
			"%"+query.Search+"%",
			"%"+query.Search+"%",
			"%"+query.Search+"%",
		)
		argIdx += 3
	}

	if query.Status != "" {
		conds = append(conds, fmt.Sprintf("status=$%d", argIdx))
		args = append(args, model.ParseStatus(query.Status))
		argIdx++
	}

	if query.Department != "" {
		conds = append(conds, fmt.Sprintf("department=$%d", argIdx))
		args = append(args, query.Department)
		argIdx++
	}

	if query.BorrowedFrom != nil {
		conds = append(conds, fmt.Sprintf("borrowed_date >= $%d", argIdx))
		args = append(args, *query.BorrowedFrom)
		argIdx++
	}

	if query.BorrowedTo != nil {
		conds = append(conds, fmt.Sprintf("borrowed_date <= $%d", argIdx))
		args = append(args, *query.BorrowedTo)
		argIdx++
	}

	if query.ReturnedFrom != nil {
		conds = append(conds, fmt.Sprintf("returned_date >= $%d", argIdx))
		args = append(args, *query.ReturnedFrom)
		argIdx++
	}

	if query.ReturnedTo != nil {
		conds = append(conds, fmt.Sprintf("returned_date <= $%d", argIdx))
		args = append(args, *query.ReturnedTo)
		argIdx++
	}

	return strings.Join(conds, " AND "), args
}

func (r *loanSlipRepository) List(ctx context.Context, query *dto.LoanSlipQuery) ([]*model.LoanSlip, error) {
	where, args := buildLoanSlipQuery(query)
	offset := (query.Page - 1) * query.Limit
	orderBy := buildOrderByClause(query.Sort, query.Order)

	sql := fmt.Sprintf(`
		SELECT id, borrower_name, department, position, name,
	    description, status, serial_number, images, borrowed_date,
		returned_date, created_by, updated_at, created_at
		FROM loan_slips
		WHERE %s
		ORDER BY %s
		LIMIT %d OFFSET %d
	`,
		where,
		orderBy,
		query.Limit,
		offset,
	)

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanLoanSlips(rows)
}

func buildOrderByClause(sort, order string) string {
	validSortColumns := map[string]bool{
		"id": true, "name": true, "borrower_name": true, "department": true,
		"status": true, "serial_number": true, "borrowed_date": true,
		"returned_date": true, "created_at": true, "updated_at": true,
	}

	if sort == "" || !validSortColumns[sort] {
		sort = "created_at"
	}

	if sort == "borrower_name" || sort == "name" || sort == "department" {
		sort = fmt.Sprintf("LOWER(%s)", sort)
	}

	return fmt.Sprintf("%s %s", sort, utils.NormalizeOrder(order))
}

func scanLoanSlips(rows interface {
	Next() bool
	Scan(...interface{}) error
	Close()
}) ([]*model.LoanSlip, error) {
	var loanSlips []*model.LoanSlip
	for rows.Next() {
		var loanSlip model.LoanSlip
		if err := rows.Scan(
			&loanSlip.ID,
			&loanSlip.BorrowerName,
			&loanSlip.Department,
			&loanSlip.Position,
			&loanSlip.Name,
			&loanSlip.Description,
			&loanSlip.Status,
			&loanSlip.SerialNumber,
			&loanSlip.Images,
			&loanSlip.BorrowedDate,
			&loanSlip.ReturnedDate,
			&loanSlip.CreatedBy,
			&loanSlip.UpdatedAt,
			&loanSlip.CreatedAt,
		); err != nil {
			return nil, err
		}
		loanSlips = append(loanSlips, &loanSlip)
	}
	return loanSlips, nil
}

func (r *loanSlipRepository) Count(ctx context.Context, query *dto.LoanSlipQuery) (int, error) {
	where, args := buildLoanSlipQuery(query)
	sql := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM loan_slips
		WHERE %s
	`, where)

	var total int
	err := r.db.QueryRow(ctx, sql, args...).Scan(&total)
	return total, err
}

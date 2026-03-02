package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/davidcm146/assets-management-be.git/internal/dto"
	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LoanSlipRepository interface {
	FindByID(ctx context.Context, id int) (*model.LoanSlip, error)
	List(ctx context.Context, query *dto.LoanSlipQuery) ([]*model.LoanSlip, error)
	Count(ctx context.Context, query *dto.LoanSlipQuery) (int, error)
	Create(ctx context.Context, loanSlip *model.LoanSlip) error
	Update(ctx context.Context, loanSlip *model.LoanSlip) error
	FindOverdue(ctx context.Context) ([]*model.LoanSlip, error)
	MarkOverdueNotified(ctx context.Context, id int) error
	Delete(ctx context.Context, id int) error
}

type loanSlipRepository struct {
	db *pgxpool.Pool
}

func NewLoanSlipRepository(db *pgxpool.Pool) LoanSlipRepository {
	return &loanSlipRepository{db: db}
}

func (r *loanSlipRepository) Create(ctx context.Context, loanSlip *model.LoanSlip) error {
	status := model.Borrowing
	_, err := r.db.Exec(ctx,
		"INSERT INTO loan_slips (name, borrower_name, department, position, description, status, serial_number, images, created_by, borrowed_date, returned_date, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		loanSlip.Name, loanSlip.BorrowerName, loanSlip.Department, loanSlip.Position, loanSlip.Description, status, loanSlip.SerialNumber, loanSlip.Images, loanSlip.CreatedBy, loanSlip.BorrowedDate, loanSlip.ReturnedDate, time.Now().UTC(), time.Now().UTC(),
	)
	return err
}

func (r *loanSlipRepository) Update(ctx context.Context, loanSlip *model.LoanSlip) error {
	var returnedAt any = loanSlip.ReturnedDate

	if loanSlip.Status == model.Returned && loanSlip.ReturnedDate == nil {
		returnedAt = time.Now().UTC()
	}

	_, err := r.db.Exec(ctx,
		`UPDATE loan_slips
		 SET name=$1,
		     borrower_name=$2,
		     department=$3,
		     position=$4,
		     description=$5,
		     status=$6,
		     serial_number=$7,
		     images=$8,
		     borrowed_date=$9,
		     returned_date=$10,
		     updated_at=$11,
			 returned_at=$12
		 WHERE id=$13`,
		loanSlip.Name,
		loanSlip.BorrowerName,
		loanSlip.Department,
		loanSlip.Position,
		loanSlip.Description,
		loanSlip.Status,
		loanSlip.SerialNumber,
		loanSlip.Images,
		loanSlip.BorrowedDate,
		loanSlip.ReturnedDate,
		time.Now().UTC(),
		returnedAt,
		loanSlip.ID,
	)

	return err
}

func (r *loanSlipRepository) Delete(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx,
		"DELETE FROM loan_slips WHERE id=$1",
		id,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("loan slip not found")
	}

	return nil
}

func (r *loanSlipRepository) FindByID(ctx context.Context, id int) (*model.LoanSlip, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, name, borrower_name, department, position, description, status, serial_number, images, created_by, borrowed_date, returned_date, updated_at, created_at
		 FROM loan_slips WHERE id = $1`, id)

	var loanSlip model.LoanSlip
	err := row.Scan(
		&loanSlip.ID,
		&loanSlip.Name,
		&loanSlip.BorrowerName,
		&loanSlip.Department,
		&loanSlip.Position,
		&loanSlip.Description,
		&loanSlip.Status,
		&loanSlip.SerialNumber,
		&loanSlip.Images,
		&loanSlip.CreatedBy,
		&loanSlip.BorrowedDate,
		&loanSlip.ReturnedDate,
		&loanSlip.UpdatedAt,
		&loanSlip.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &loanSlip, nil
}

func (r *loanSlipRepository) FindOverdue(ctx context.Context) ([]*model.LoanSlip, error) {
	query := `
		SELECT id, borrower_name, department, position,
		       name, description, status, serial_number,
		       images, borrowed_date, returned_date,
		       created_by, updated_at, created_at
		FROM loan_slips
		WHERE DATE(returned_date) < CURRENT_DATE
		AND status = $1
		AND overdue_notified = FALSE
	`

	rows, err := r.db.Query(ctx, query, model.Borrowing)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.LoanSlip

	for rows.Next() {
		var loan model.LoanSlip
		err := rows.Scan(
			&loan.ID,
			&loan.BorrowerName,
			&loan.Department,
			&loan.Position,
			&loan.Name,
			&loan.Description,
			&loan.Status,
			&loan.SerialNumber,
			&loan.Images,
			&loan.BorrowedDate,
			&loan.ReturnedDate,
			&loan.CreatedBy,
			&loan.UpdatedAt,
			&loan.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &loan)
	}

	return results, nil
}

func (r *loanSlipRepository) MarkOverdueNotified(ctx context.Context, id int) error {
	query := `
		UPDATE loan_slips
		SET overdue_notified = TRUE,
		    updated_at = $2, overdue_at = $3
		WHERE id = $1
		  AND overdue_notified = FALSE
	`
	cmdTag, err := r.db.Exec(ctx, query, id, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return nil
	}

	return nil
}

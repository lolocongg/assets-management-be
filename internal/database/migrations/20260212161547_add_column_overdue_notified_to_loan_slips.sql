-- +goose Up
-- +goose StatementBegin
ALTER TABLE loan_slips
ADD COLUMN overdue_notified BOOLEAN DEFAULT FALSE;

CREATE INDEX idx_loan_slips_overdue_check
ON loan_slips (returned_date, status, overdue_notified);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_loan_slips_overdue_check;

ALTER TABLE loan_slips
DROP COLUMN IF EXISTS overdue_notified;
-- +goose StatementEnd

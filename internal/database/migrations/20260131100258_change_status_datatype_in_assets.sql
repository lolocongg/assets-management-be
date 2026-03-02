-- +goose Up
ALTER TABLE loan_slips
ALTER COLUMN status TYPE integer USING status::integer;

-- +goose Down
ALTER TABLE loan_slips
ALTER COLUMN status TYPE VARCHAR USING status::VARCHAR;

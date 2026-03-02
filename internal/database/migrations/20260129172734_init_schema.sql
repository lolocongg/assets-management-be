-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    role INTEGER NOT NULL,
    created_at TIMESTAMP
);

CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    action VARCHAR,
    device_name VARCHAR,
    ip_address VARCHAR,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE loan_slips (
    id SERIAL PRIMARY KEY,
    borrower_name VARCHAR,
    department VARCHAR,
    position VARCHAR,
    name VARCHAR,
    description TEXT,
    status VARCHAR,
    serial_number VARCHAR,
    images TEXT[],
    borrowed_date TIMESTAMP WITH TIME ZONE,
    returned_date TIMESTAMP WITH TIME ZONE,
    created_by INTEGER NOT NULL,
    returned_at TIMESTAMP WITH TIME ZONE,
    overdue_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    recipient_id INTEGER,
    sender_id INTEGER,
    title VARCHAR,
    type INTEGER,
    payload JSONB,
    content VARCHAR,
    is_read BOOLEAN,
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Foreign key constraints
ALTER TABLE loan_slips
    ADD CONSTRAINT fk_loan_slips_created_by
    FOREIGN KEY (created_by) REFERENCES users(id);

ALTER TABLE audit_logs
    ADD CONSTRAINT fk_audit_logs_user
    FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE notifications
    ADD CONSTRAINT fk_notifications_recipient
    FOREIGN KEY (recipient_id) REFERENCES users(id);

ALTER TABLE notifications
    ADD CONSTRAINT fk_notifications_sender
    FOREIGN KEY (sender_id) REFERENCES users(id);

CREATE INDEX idx_loan_slips_created_by
ON loan_slips (created_by);

CREATE INDEX idx_loan_slips_position
ON loan_slips (position);

CREATE INDEX idx_loan_slips_department
ON loan_slips (department);

CREATE INDEX idx_loan_slips_borrower_name
ON loan_slips (borrower_name);

CREATE INDEX idx_loan_slips_borrowed_date
ON loan_slips (borrowed_date DESC);

CREATE INDEX idx_loan_slips_created_at
ON loan_slips (created_at DESC);

CREATE INDEX idx_notifications_payload
ON notifications USING GIN (payload);

-- +goose Down
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS loan_slips;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS users;
DROP INDEX IF EXISTS idx_loan_slips_borrowed_date;
DROP INDEX IF EXISTS idx_loan_slips_created_at;
DROP INDEX IF EXISTS idx_loan_slips_created_by;
DROP INDEX IF EXISTS idx_loan_slips_position;
DROP INDEX IF EXISTS idx_loan_slips_department;
DROP INDEX IF EXISTS idx_loan_slips_borrower_name;
DROP INDEX IF EXISTS idx_notifications_payload;

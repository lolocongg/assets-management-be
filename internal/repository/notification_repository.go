package repository

import (
	"context"
	// "encoding/json"
	"fmt"
	"time"

	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository interface {
	Create(ctx context.Context, n *model.Notification) error
	ListByRecipient(ctx context.Context, recipientID int, offset int, limit int, isRead *bool) ([]*model.Notification, int, error)
	MarkAsRead(ctx context.Context, id int) error
	CountUnread(ctx context.Context, recipientID int) (int, error)
	BulkCreate(ctx context.Context, notifications []*model.Notification) error
}

type notificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, n *model.Notification) error {
	query := `
		INSERT INTO notifications
		(recipient_id, sender_id, title, type, payload, content, is_read, created_at)
		VALUES ($1,$2,$3,$4,$5,false,$6)
		RETURNING id
	`

	return r.db.QueryRow(ctx, query,
		n.RecipientID,
		n.SenderID,
		n.Title,
		n.Type,
		n.Payload,
		n.Content,
		time.Now().UTC(),
	).Scan(&n.ID)
}

func (r *notificationRepository) ListByRecipient(ctx context.Context, recipientID int, offset int, limit int, isRead *bool) ([]*model.Notification, int, error) {
	baseQuery := `
		FROM notifications
		WHERE recipient_id = $1
	`
	args := []interface{}{recipientID}
	argIndex := 2

	if isRead != nil {
		baseQuery += fmt.Sprintf(" AND is_read = $%d", argIndex)
		args = append(args, *isRead)
		argIndex++
	}

	countQuery := "SELECT COUNT(*) " + baseQuery

	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	listQuery := `
		SELECT id, recipient_id, sender_id, title, type,
		       content, payload, is_read, read_at, created_at
	` + baseQuery +
		fmt.Sprintf(" ORDER BY created_at DESC OFFSET $%d LIMIT $%d", argIndex, argIndex+1)

	args = append(args, offset, limit)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	notifications := make([]*model.Notification, 0)

	for rows.Next() {
		var n model.Notification

		if err := rows.Scan(
			&n.ID,
			&n.RecipientID,
			&n.SenderID,
			&n.Title,
			&n.Type,
			&n.Content,
			&n.Payload,
			&n.IsRead,
			&n.ReadAt,
			&n.CreatedAt,
		); err != nil {
			return nil, 0, err
		}

		notifications = append(notifications, &n)
	}

	if notifications == nil {
		notifications = []*model.Notification{}
	}

	return notifications, total, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id int) error {
	query := `
		UPDATE notifications
		SET is_read = true, read_at = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}

func (r *notificationRepository) CountUnread(ctx context.Context, recipientID int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE recipient_id = $1 AND is_read = false
	`

	var count int
	err := r.db.QueryRow(ctx, query, recipientID).Scan(&count)
	return count, err
}

func (r *notificationRepository) BulkCreate(ctx context.Context, notifications []*model.Notification) error {

	if len(notifications) == 0 {
		return nil
	}

	now := time.Now().UTC()
	rows := make([][]interface{}, 0, len(notifications))

	for _, n := range notifications {
		// payload := n.Payload
		// if payload == nil {
		// 	payload = json.RawMessage(`{}`)
		// }
		rows = append(rows, []interface{}{
			n.RecipientID,
			n.SenderID,
			n.Title,
			n.Type,
			n.Content,
			string(n.Payload),
			false,
			now,
		})
	}

	_, err := r.db.CopyFrom(
		ctx,
		pgx.Identifier{"notifications"},
		[]string{
			"recipient_id",
			"sender_id",
			"title",
			"type",
			"content",
			"payload",
			"is_read",
			"created_at",
		},
		pgx.CopyFromRows(rows),
	)

	return err
}

package dto

type CreateNotificationRequest struct {
	RecipientID int    `json:"recipient_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Type        int    `json:"type"`
	Content     string `json:"content"`
}

type NotificationQuery struct {
	Page   int   `form:"page"`
	Limit  int   `form:"limit"`
	IsRead *bool `form:"is_read"`
}

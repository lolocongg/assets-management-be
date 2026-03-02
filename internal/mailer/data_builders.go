package email

import (
	"encoding/json"

	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/davidcm146/assets-management-be.git/internal/utils"
)

type OverdueItem struct {
	BorrowerName string
	BorrowedDate string
	ReturnedDate string
}

type OverdueEmailData struct {
	Total int
	More  int
	Items []OverdueItem
}

func BuildOverdueEmailData(notifications []*model.Notification, limit int) OverdueEmailData {
	total := len(notifications)

	display := notifications
	more := 0

	if total > limit {
		display = notifications[:limit]
		more = total - limit
	}

	items := make([]OverdueItem, 0, len(display))
	for _, n := range display {
		var payload model.NotificationPayload

		if err := json.Unmarshal(n.Payload, &payload); err != nil {
			continue
		}
		borrowerName, _ := payload.Extra["borrower_name"].(string)
		returnedRaw, _ := payload.Extra["returned_date"].(string)
		borrowedRaw, _ := payload.Extra["borrowed_date"].(string)

		items = append(items, OverdueItem{
			BorrowerName: borrowerName,
			BorrowedDate: utils.FormatDate(borrowedRaw),
			ReturnedDate: utils.FormatDate(returnedRaw),
		})
	}

	return OverdueEmailData{
		Total: total,
		More:  more,
		Items: items,
	}
}

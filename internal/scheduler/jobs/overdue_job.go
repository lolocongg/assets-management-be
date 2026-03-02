package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/davidcm146/assets-management-be.git/internal/service"
)

type OverdueJob struct {
	loanSlipService     service.LoanSlipService
	notificationService service.NotificationService
}

func NewOverdueJob(loanSlipService service.LoanSlipService, notificationService service.NotificationService) *OverdueJob {
	return &OverdueJob{
		loanSlipService:     loanSlipService,
		notificationService: notificationService,
	}
}

func (j *OverdueJob) Name() string {
	return "overdue-loan-job"
}

func (j *OverdueJob) Schedule() string {
	return "0 9 * * *"
}

func (j *OverdueJob) Run() {
	log.Println("[JOB START]", j.Name())
	ctx := context.Background()

	overdues, err := j.loanSlipService.GetOverdue(ctx)
	if err != nil {
		log.Println("[JOB ERROR] get overdue:", err)
		return
	}

	if len(overdues) == 0 {
		log.Println("[JOB DONE] no overdue loan slips")
		return
	}

	notifications := make([]*model.Notification, 0, len(overdues))
	for _, slip := range overdues {
		_, err := j.loanSlipService.MarkAsOverdue(ctx, slip.ID)
		if err != nil {
			continue
		}

		payloadObj := model.NotificationPayload{
			Entity: "loan_slip",
			Action: "navigate",
			URL:    fmt.Sprintf("/loan-slips/%d", slip.ID),
			Extra: map[string]interface{}{
				"id":            slip.ID,
				"borrower_name": slip.BorrowerName,
				"borrowed_date": slip.BorrowedDate,
				"returned_date": slip.ReturnedDate,
			},
		}
		payloadBytes, err := json.Marshal(payloadObj)
		if err != nil {
			log.Println("marshal payload failed:", err)
			continue
		}

		notifications = append(notifications, &model.Notification{
			RecipientID: slip.CreatedBy,
			SenderID:    nil,
			Title:       "Phiếu mượn quá hạn",
			Type:        int(model.NotificationTypeLoanSlipOverdue),
			Content:     "Phiếu mượn tài sản \"" + slip.Name + "\" đã quá hạn trả.",
			Payload:     payloadBytes,
		})
	}

	_, err = j.notificationService.BulkSend(ctx, notifications)
	if err != nil {
		log.Println(err)
	}

	for _, slip := range overdues {
		j.loanSlipService.MarkOverdueNotified(ctx, slip.ID)
	}

	log.Printf("[JOB DONE] processed %d overdue loan slips\n", len(overdues))
}

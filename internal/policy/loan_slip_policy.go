package policy

import "github.com/davidcm146/assets-management-be.git/internal/model"

var adminUpdateFields = map[string]bool{
	"status":        true,
	"borrowed_date": true,
	"returned_date": true,
}

var itUpdateFields = map[string]bool{
	"owner_name":    true,
	"department":    true,
	"position":      true,
	"name":          true,
	"description":   true,
	"status":        true,
	"serial_number": true,
	"images":        true,
	"borrowed_date": true,
	"returned_date": true,
}

type LoanSlipPolicy struct{}

func NewLoanSlipPolicy() *LoanSlipPolicy {
	return &LoanSlipPolicy{}
}

func (p *LoanSlipPolicy) allowedFields(role string) map[string]bool {
	switch role {
	case model.Admin.String():
		return adminUpdateFields
	case model.IT.String():
		return itUpdateFields
	default:
		return nil
	}
}

func (p *LoanSlipPolicy) ForbiddenFields(role string, fields []string) []string {
	allowed := p.allowedFields(role)
	if allowed == nil {
		return fields // role không hợp lệ → cấm hết
	}

	var forbidden []string
	for _, field := range fields {
		if !allowed[field] {
			forbidden = append(forbidden, field)
		}
	}
	return forbidden
}

package model

import "time"

type TimeFilter struct {
	From *time.Time `json:"from"`
	To   *time.Time `json:"to"`
}

type LoanMetrics struct {
	Total     int64 `json:"total"`
	Borrowing int64 `json:"borrowing"`
	Returned  int64 `json:"returned"`
	Overdue   int64 `json:"overdue"`
}

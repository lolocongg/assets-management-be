package dto

type DashboardFilterRequest struct {
	From   string `form:"from"   json:"from"`
	To     string `form:"to"     json:"to"`
	Period string `form:"period" json:"period"`
}

type LoanMetricsResponse struct {
	Total     int64 `json:"total"`
	Borrowing int64 `json:"borrowing"`
	Returned  int64 `json:"returned"`
	Overdue   int64 `json:"overdue"`
}

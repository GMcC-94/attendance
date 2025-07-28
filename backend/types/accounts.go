package types

import "time"

type AccountEntry struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"` // income or expenditure
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateAccountsRequest struct {
	Income      []AccountEntry `json:"income"`
	Expenditure []AccountEntry `json:"expenditure"`
}

type CreateAccountsResponse struct {
	Message string `json:"message"`
}

type GroupedAccounts struct {
	Year   int                    `json:"year"`
	Months []GroupedMonthAccounts `json:"months"`
}

type GroupedMonthAccounts struct {
	Month   string               `json:"month"`
	Entries []GroupedDayAccounts `json:"entries"`
}

type GroupedDayAccounts struct {
	Date    string         `json:"date"`
	Records []AccountEntry `json:"records"`
}

package models

type Wallet struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type Category struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Type   string `json:"type"` // "INCOME" atau "EXPENSE"
}
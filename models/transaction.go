package models

type Transaction struct {
	ID       int    `json:"id"`
	UserID   string `json:"-"`
	Title    string `json:"title"`
	Amount   int    `json:"amount"`
	Type     string `json:"type"`
	DueDate  string `json:"due_date"`
	GroupID  string `json:"group_id"`
	Tenor    int    `json:"tenor"`    
	Status   string `json:"status"` // TAMBAHIN INI: "PENDING" atau "PAID"
}
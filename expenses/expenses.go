package expenses

import "database/sql"

var db *sql.DB
var tableName string

func SetDB(d *sql.DB) {
	db = d
}

func SetTableName(t string) {
	tableName = t
}

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type ExpenseResponse struct {
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

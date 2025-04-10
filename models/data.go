package models

type Data struct {
	IsAuth   bool
	Error    string
	Expenses Expenses
	User
}

package models

type Expenses struct {
	TodayExpenses []*Expense
	TodayTotal    float64
	MonthExpenses []*Expense
	MonthTotal    float64
}

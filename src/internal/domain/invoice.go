package domain

import "time"

type Invoice struct {
	ID              int64
	CustomerEmail   string
	CustomerName    string
	CustomerAddress string
	CustomerId      string
	TotalAmount     float64
	Subtotal        float64
	TaxAmount       float64
	Items           []*InvoiceItem
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DueDate         time.Time
	Month           int
	Year            int
	Sent            bool
}
type InvoiceItem struct {
	ID          int64
	InvoiceID   int64
	Description string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}

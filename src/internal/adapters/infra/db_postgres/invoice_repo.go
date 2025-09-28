package db_postgres

import (
	"database/sql"

	"github.com/joseMChavez/fc-job/src/internal/domain"
	_ "github.com/lib/pq"
)

type InvoiceRepo struct {
	db *sql.DB
}

func NewInvoiceRepo(db *sql.DB) *InvoiceRepo {
	return &InvoiceRepo{db: db}
}

func (r *InvoiceRepo) FindPending() ([]*domain.Invoice, error) {
	rows, err := r.db.Query("SELECT * FROM invoices WHERE sent = false")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*domain.Invoice
	for rows.Next() {
		var inv *domain.Invoice
		if err := rows.Scan(&inv.ID, &inv.CustomerEmail, &inv.CustomerName, &inv.CustomerAddress, &inv.CustomerId, &inv.CreatedAt,
			&inv.Month, &inv.Year, &inv.Subtotal, &inv.TotalAmount, &inv.TaxAmount, &inv.Sent, &inv.DueDate); err != nil {
			return nil, err
		}
		res = append(res, inv)
	}
	return res, nil
}

func (r *InvoiceRepo) MarkAsSent(id int64) error {
	_, err := r.db.Exec("UPDATE invoices SET sent = true, sent_at = NOW() WHERE id = $1", id)
	return err
}
func (r *InvoiceRepo) GetDetails(id int64) ([]*domain.InvoiceItem, error) {
	rows, err := r.db.Query("SELECT id, invoice_id, description, quantity, unit_price, total_price FROM invoice_items WHERE invoice_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*domain.InvoiceItem
	for rows.Next() {
		var item *domain.InvoiceItem
		if err := rows.Scan(&item.ID, &item.InvoiceID, &item.Description, &item.Quantity, &item.UnitPrice, &item.TotalPrice); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
func (r *InvoiceRepo) GetByID(id int64) (*domain.Invoice, error) {
	var inv domain.Invoice
	err := r.db.QueryRow("SELECT * FROM invoices WHERE id = $1", id).
		Scan(&inv.ID, &inv.CustomerEmail, &inv.CustomerName, &inv.CustomerAddress, &inv.CustomerId, &inv.CreatedAt, &inv.DueDate,
			&inv.Month, &inv.Year, &inv.Subtotal, &inv.TotalAmount, &inv.TaxAmount, &inv.Sent)
	if err != nil {

		return nil, err
	}
	inv.Items, err = r.GetDetails(inv.ID)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

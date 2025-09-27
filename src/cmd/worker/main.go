package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joseMChavez/fc-job/internal/adapters/infra/db_postgres"
	"github.com/joseMChavez/fc-job/internal/adapters/infra/pdf_gofpdf"
	"github.com/joseMChavez/fc-job/internal/adapters/infra/smtp_outlook"
	"github.com/joseMChavez/fc-job/internal/application"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := db_postgres.NewInvoiceRepo(db)
	sender := smtp_outlook.NewOutlookSender(
		os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"),
	)
	pdf := pdf_gofpdf.NewPDFGenerator()

	uc := application.InvoiceSender{
		Repo:        repo,
		EmailSender: sender,
		PdfGen:      pdf,
	}

	if err := uc.SendInvoice(); err != nil {
		log.Fatal(err)
	}
}

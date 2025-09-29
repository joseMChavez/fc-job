package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joseMChavez/fc-job/src/internal/adapters/infra/db_postgres"
	"github.com/joseMChavez/fc-job/src/internal/adapters/infra/pdf_gofpdf"
	"github.com/joseMChavez/fc-job/src/internal/adapters/infra/smtp_outlook"
	"github.com/joseMChavez/fc-job/src/internal/application"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
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
	cc := cron.New(cron.WithSeconds())

	_, err = cc.AddFunc("0 */2 * * * *", func() {

		if err := uc.SendInvoice(); err != nil {
			log.Fatal(err)
			log.Println("Error enviando facturas:", err)
		}

	})
	if err != nil {
		log.Fatal(err)
	}
	cc.Start()

	fmt.Println("Cron iniciado. Presiona Ctrl+C para salir.")
	select {} // ejecucion indefinida
}

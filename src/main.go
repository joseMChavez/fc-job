package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

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

	//repo := db_postgres.NewInvoiceRepo(db) 	"github.com/joseMChavez/fc-job/src/internal/adapters/infra/db_postgres"

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Error al conectar a DB: %v", err)
	}
	defer pool.Close()

	// Crear tabla si no existe
	_, err = pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS cron_logs (
            id SERIAL PRIMARY KEY,
            message TEXT,
            created_at TIMESTAMP DEFAULT NOW()
        )
    `)
	if err != nil {
		log.Fatalf("Error al crear tabla: %v", err)
	}
	//sender := smtp_outlook.NewOutlookSender(
	//	os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"),
	//	os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"),
	//)
	/*pdf := pdf_gofpdf.NewPDFGenerator()

	uc := application.InvoiceSender{
		Repo:        repo,
		EmailSender: sender,
		PdfGen:      pdf,
	}*/
	cc := cron.New(cron.WithSeconds())

	_, err = cc.AddFunc("0 */2 * * * *", func() {
		message := fmt.Sprintf("Cron ejecutado a las %s", time.Now().Format(time.RFC3339))
		fmt.Println(message)

		_, err := pool.Exec(context.Background(),
			"INSERT INTO cron_logs(message) VALUES($1)", message)
		if err != nil {
			log.Printf("Error al guardar en DB: %v", err)
		}
		/*if err := uc.SendInvoice(); err != nil {
			log.Fatal(err)
			log.Println("Error enviando facturas:", err)
		}*/

	})
	if err != nil {
		log.Fatal(err)
	}
	cc.Start()

	fmt.Println("Cron iniciado. Presiona Ctrl+C para salir.")
	select {} // ejecucion indefinida
}

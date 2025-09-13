package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"decentraton/internal"
	"decentraton/internal/models"
	"decentraton/pkg/config"
	"decentraton/pkg/database"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := internal.RunMigrations(cfg.PGURL); err != nil {
		log.Fatal(err)
	}

	data, err := internal.ReadCsvFile("case 1")
	if err != nil {
		log.Fatal(err)
	}

	if err := insert(db, data); err != nil {
		log.Fatal(err)
	}

	allAnalysis, err := internal.GetAllClientAnalysis(db)
	if err != nil {
		log.Fatal(err)
	}

	for clientCode, analysis := range allAnalysis {
		fmt.Printf("Client %d: %+v\n", clientCode, analysis)
	}
}

func insert(db *sql.DB, data []models.CSVData) error {
	var clientCode int
	if err := db.QueryRow(`SELECT client_code FROM clients LIMIT 1`).Scan(&clientCode); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("could not select client code: %w", err)
	}
	if clientCode != 0 {
		fmt.Println("DATABASE already full filled")
		return nil
	}

	for _, d := range data {
		switch {
		case strings.Contains(d.FileName, "clients"):
			if err := insertClients(db, d.Data); err != nil {
				return err
			}
		case strings.Contains(d.FileName, "transfers"):
			if err := insertTransfers(db, d.Data); err != nil {
				return err
			}

		case strings.Contains(d.FileName, "transactions"):
			if err := insertTransaction(db, d.Data); err != nil {
				return err
			}
		}
	}

	return nil
}

func insertClients(db *sql.DB, data [][]string) error {
	for _, row := range data {
		_, err := db.Exec(`INSERT INTO clients 
    (client_code, name, status, age, city, avg_monthly_balance_KZT) 
VALUES
($1, $2, $3, $4, $5, $6);`,
			row[0], row[1], row[2], row[3], row[4], row[5])
		if err != nil {
			return fmt.Errorf("insertClients: %w", err)
		}
	}

	return nil
}

func insertTransfers(db *sql.DB, data [][]string) error {
	for _, row := range data {
		_, err := db.Exec(`
INSERT INTO transfers 
(client_code, name, product, status, city, date, type, direction, amount, currency) 
VALUES
($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
			row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], row[9])
		if err != nil {
			return fmt.Errorf("insertTransfers: %w", err)
		}
	}

	return nil
}

func insertTransaction(db *sql.DB, data [][]string) error {
	for _, row := range data {
		_, err := db.Exec(`
INSERT INTO transactions
    (client_code, name, product, status, city, date, category, amount, currency) 
VALUES
($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
			row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8])
		if err != nil {
			return fmt.Errorf("insertTransactions: %w", err)
		}
	}

	return nil
}

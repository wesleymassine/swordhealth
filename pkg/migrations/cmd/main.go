package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wesleymassine/swordhealth/migrations"
)
func main() {
    db, err := NewMySQLConnection()
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }
    defer db.Close()

    log.Println("Applying migrations...")
    err = migrations.ApplyMigrations(db)
    if err != nil {
        log.Fatalf("Could not apply migrations: %v", err)
    }
    log.Println("Migrations applied successfully.")
}

func NewMySQLConnection() (*sql.DB, error) {
    db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/task")
    if err != nil {
        return nil, err
    }
    return db, nil
}
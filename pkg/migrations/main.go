package migrations

import (
    "database/sql"
    "log"

    // _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/task")
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }
    defer db.Close()

    log.Println("Applying migrations...")
    err = ApplyMigrations(db)
    if err != nil {
        log.Fatalf("Could not apply migrations: %v", err)
    }
    log.Println("Migrations applied successfully.")
}

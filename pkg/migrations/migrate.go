package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrationsUp applies all SQL migration files using golang-migrate
func RunMigrationsUp() {
	// Create the database schema if it doesn't exist
	createDatabase("root:root@tcp(localhost:3306)/")

	// Set the DB URL
	dbURL := "mysql://root:root@tcp(localhost:3306)/task"

	if dbURL == "" {
		log.Fatalf("DB_URL is not set")
	}

	// Path to the migrations folder
	migrationsDir := "file://./sql" // file:// for local files

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get working directory: %v", err)
	}
	migrationsPath := filepath.Join(workingDir, migrationsDir)
	fmt.Println("Migrations Path:", migrationsPath)

	// Create new migration instance
	m, err := migrate.New(migrationsDir, dbURL)
	if err != nil {
		log.Fatalf("Could not create new migrate instance: %v", err)
	}

	// Run migrations up
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not apply migrations: %v", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}

	// Close the migrate instance to ensure proper cleanup
	m.Close()
}

// RunMigrationsDown rolls back the most recent migration using golang-migrate
func RunMigrationsDown() {
	// Set the DB URL
	dbURL := "mysql://root:root@tcp(localhost:3306)/task"

	if dbURL == "" {
		log.Fatalf("DB_URL is not set")
	}

	// Path to the migrations folder
	migrationsDir := "file://./sql" // file:// for local files

	// Create new migration instance
	m, err := migrate.New(migrationsDir, dbURL)
	if err != nil {
		log.Fatalf("Could not create new migrate instance: %v", err)
	}

	// Roll back the most recent migration (down)
	err = m.Steps(-1) // Rollback one step
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not roll back migrations: %v", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No migrations to roll back")
	} else {
		log.Println("Migrations rolled back successfully")
	}

	// Close the migrate instance to ensure proper cleanup
	m.Close()
}

// createDatabase checks if the 'task' database exists, and creates it if it doesn't
func createDatabase(dsn string) {
	// Connect to MySQL server (without specifying a database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL server: %v", err)
	}
	defer db.Close()

	// Check if 'task' database exists, and create it if it doesn't
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS task")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	log.Println("Database 'task' exists or has been created.")
}

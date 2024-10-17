package migrations

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
)

const migrationsDir = "./sql/"

// NewMigrations applies all migrations
func NewMigrations() {
	db, err := NewMySQLConnection() // Establish DB connection
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	log.Println("Applying migrations...")
	err = ApplyMigrations(db)

	if err != nil {
		log.Fatalf("Could not apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully.")
}

// ApplyMigrations applies all SQL migration files in sequence
func ApplyMigrations(db *sql.DB) error {
	// Get the current working directory and create the full path to migrations
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working directory: %v", err)
	}
	migrationsPath := filepath.Join(workingDir, migrationsDir)

	fmt.Println("DEBUG", migrationsPath)

	// Read all migration files
	files, err := readMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("could not read migrations directory: %v", err)
	}

	// Execute all .sql files
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			err = applyMigration(db, filepath.Join(migrationsPath, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// readMigrationFiles reads all the migration files in the given directory
func readMigrationFiles(dir string) ([]fs.FileInfo, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	list, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	// Sort files by name (i.e., run them in order)
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name() < list[j].Name()
	})
	return list, nil
}

// applyMigration executes a migration file
func applyMigration(db *sql.DB, filePath string) error {
	query, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read migration file %s: %v", filePath, err)
	}

	_, err = db.Exec(string(query))
	if err != nil {
		return fmt.Errorf("could not apply migration %s: %v", filePath, err)
	}
	log.Printf("Applied migration: %s", filePath)
	return nil
}

// NewMySQLConnection creates a new MySQL connection
func NewMySQLConnection() (*sql.DB, error) {
	dsn := os.Getenv("DB_DSN")
	return sql.Open("mysql", dsn)
}

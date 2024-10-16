package migrations

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"io"
	"sort"
)

const migrationsDir = "./pkg/migrations/sql"

func NewMigrations(db *sql.DB) {
	log.Println("Applying migrations...")
	err := ApplyMigrations(db)

	if err != nil {
		log.Fatalf("Could not apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully.")
}

// ApplyMigrations applies all SQL migration files in sequence
func ApplyMigrations(db *sql.DB) error {
	files, err := func() ([]fs.FileInfo, error) {
		f, err := os.Open(migrationsDir)
		if err != nil {
			return nil, err
		}
		list, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			return nil, err
		}
		sort.Slice(list, func(i, j int) bool {
			return list[i].Name() < list[j].Name()
		})
		return list, nil
	}()
	if err != nil {
		return fmt.Errorf("could not read migrations directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			err = applyMigration(db, filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// applyMigration runs a single SQL migration file
func applyMigration(db *sql.DB, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("could not open migration file %s: %v", filepath, err)
	}
	defer file.Close()

	sqlBytes, err := io.ReadAll(io.Reader(file))
	if err != nil {
		return fmt.Errorf("could not read migration file %s: %v", filepath, err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("could not execute migration file %s: %v", filepath, err)
	}

	log.Printf("Migration applied: %s", filepath)
	return nil
}

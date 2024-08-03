package migrate

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

// Migration represents a database migration
type Migration struct {
	ID   string
	Up   string
	Down string
}

// LoadMigrations loads migration files from a directory
func LoadMigrations(dir string) ([]Migration, error) {
	var migrations []Migration
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			content, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to read migration file: %w", err)
			}
			migrations = append(migrations, Migration{
				ID: file.Name(),
				Up: string(content),
			})
		}
	}
	return migrations, nil
}

// RunMigrations runs all migrations on the database
func RunMigrations(db *sql.DB, migrations []Migration) error {
	// Ensure the migrations table exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
        id VARCHAR(255) PRIMARY KEY,
        applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	for _, migration := range migrations {
		// Check if the migration has already been applied
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM migrations WHERE id = ?)", migration.ID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if migration %s has been applied: %w", migration.ID, err)
		}
		if exists {
			log.Printf("Migration %s already applied, skipping", migration.ID)

			continue
		}

		// Apply the migration
		_, err = db.Exec(migration.Up)
		if err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migration.ID, err)
		}

		// Record the migration as applied
		_, err = db.Exec("INSERT INTO migrations (id) VALUES (?)", migration.ID)
		if err != nil {
			return fmt.Errorf("failed to record migration %s as applied: %w", migration.ID, err)
		}

		log.Printf("Migration %s applied successfully", migration.ID)
	}
	return nil
}

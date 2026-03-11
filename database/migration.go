package database

import (
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// RunDBMigration - Chạy tất cả pending migrations từ embedded SQL files
func RunDBMigration(dbSource string) {
	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		log.Fatalf("[Migration] Failed to load embedded migration files: %v", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dbSource)
	if err != nil {
		log.Fatalf("[Migration] Failed to create migrate instance: %v", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("[Migration] Source close error: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("[Migration] DB close error: %v", dbErr)
		}
	}()

	currentVersion, dirty, _ := m.Version()
	if dirty {
		log.Fatalf("[Migration] DB is in dirty state at version %d. Fix: make migrate-force V=%d", currentVersion, currentVersion)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Printf("[Migration] No new migrations (current version: %d)", currentVersion)
			return
		}
		log.Fatalf("[Migration] Failed to run migrations: %v", err)
	}

	newVersion, _, _ := m.Version()
	log.Printf("[Migration] ✓ Migrated from version %d → %d", currentVersion, newVersion)
}

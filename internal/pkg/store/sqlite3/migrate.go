package sqlite3

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3" // allow sqlite3 as a migration destination
	_ "github.com/golang-migrate/migrate/v4/source/file"      // allow file:// as a migration source
)

// Migrate enacts all SQL migration files
func Migrate() error {
	db, err := Connection()
	if err != nil {
		return err
	}

	driver, err := sqlite3.WithInstance(db.DB(), &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "sqlite3", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

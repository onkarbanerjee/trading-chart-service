package migrations

import (
	"embed"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"gorm.io/gorm"
)

var (
	//go:embed scripts
	migrations embed.FS
)

func New(db *gorm.DB) (*migrate.Migrate, error) {
	dbDriver, err := getDbDriver(db)
	if err != nil {
		return nil, err
	}

	source, err := httpfs.New(http.FS(migrations), "scripts")
	if err != nil {
		return nil, err
	}

	return migrate.NewWithInstance("httpfs", source, "postgres", dbDriver)
}

func getDbDriver(db *gorm.DB) (database.Driver, error) {
	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	return postgres.WithInstance(sqlDb, &postgres.Config{})
}

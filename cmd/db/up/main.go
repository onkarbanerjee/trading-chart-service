package main

import (
	"log"

	"github.com/onkarbanerjee/trading-chart-service/internal/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB() (*gorm.DB, error) {
	// Replace with your actual PostgreSQL credentials
	dsn := "host=localhost user=dbuser password=dbpass dbname=trading-chart-service port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	db, err := GetDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	migrate, err := migrations.New(db)
	if err != nil {
		log.Fatalf("error while getting migrations, %s", err.Error())
	}
	log.Println("connected, will perform migration")
	err = migrate.Up()
	if err != nil {
		log.Fatalf("error while performing migration, %s", err.Error())
	}

	log.Println("migration completed successfully")
}

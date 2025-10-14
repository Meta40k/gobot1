// db.go
package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openDB() (*gorm.DB, error) {
	dsn := os.Getenv("POSTGRES_DSN")
	fmt.Println(dsn)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&AdminlogEvent{})
}

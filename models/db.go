package services

import (
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
)

// ConnectDB establishes a connection to the PostgreSQL database using GORM.
func ConnectDB() (*gorm.DB, error) {
  dsn := "host=localhost user=postgres password=dbpass dbname=gca port=5433 sslmode=disable TimeZone=Europe/Warsaw"
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
    return nil, err
  }
  return db, nil
}

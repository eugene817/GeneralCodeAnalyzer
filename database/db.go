package database

import (
  "fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
  "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// InitDB establishes a connection to the PostgreSQL database using GORM.
// And initializes admin
func InitDB() (*gorm.DB, error) {
  host := os.Getenv("DB_HOST")
  port := os.Getenv("DB_PORT")
  user := os.Getenv("DB_USER")
  pass := os.Getenv("DB_PASSWORD")
  name := os.Getenv("DB_NAME")
  tz   := os.Getenv("DB_TIMEZONE")

  if host == "" || port == "" || user == "" || name == "" {
    return nil, fmt.Errorf("DB_HOST, DB_PORT, DB_USER and DB_NAME must be set")
  }

  dsn := fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
    host, port, user, pass, name, tz,
    )

  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
    return nil, err
  }


  if err := db.AutoMigrate(&User{}); err != nil {
    log.Fatalf("AutoMigrate user: %v", err)
  }

  adminUser := os.Getenv("ADMIN_USER")
  adminPass := os.Getenv("ADMIN_PASS")

  if adminUser == "" { adminUser = "admin" }
  if adminPass == "" { adminPass = "secret" }

  hash, err := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
  if err != nil {
    log.Fatalf("bcrypt error: %v", err)
  }

  var u User
  db.
    Where(User{Username: adminUser}).    
    Attrs(User{PasswordHash: string(hash)}).
    FirstOrCreate(&u)

  return db, nil
}

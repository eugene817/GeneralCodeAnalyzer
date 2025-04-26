package database

import (
  "gorm.io/gorm"
)

type User struct {
  gorm.Model
  Username string `gorm:"uniqueIndex;not null"`
  PasswordHash string `gorm:"not null"`
}

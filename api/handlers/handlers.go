package handlers

import (
	"github.com/eugene817/GeneralCodeAnalyzer/services"
	"gorm.io/gorm"
)

type Handler struct {
  svc *services.Service 
  db *gorm.DB
}

func NewHandler(svc *services.Service, db *gorm.DB) *Handler {
  return &Handler{
    svc: svc,
    db: db,
  }
}

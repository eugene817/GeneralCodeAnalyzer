package handlers

import (
  "github.com/eugene817/GeneralCodeAnalyzer/services"
)

type Handler struct {
  svc *services.Service 
}

func NewHandler(svc *services.Service) *Handler {
  return &Handler{
    svc: svc,
  }
}

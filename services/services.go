package services

import "github.com/eugene817/Cowdocs/api"

type Service struct {
  apiSvc *api.API
}

func NewService(apiSvc *api.API) *Service {
  return &Service{
    apiSvc: apiSvc,
  }
}

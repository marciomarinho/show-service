package service

import "github.com/marciomarinho/show-service/internal/model"

type Service interface {
	CreateShows(request model.Request) error
}

type showService struct {
}

func NewShowService() Service {
	return &showService{}
}

func (s *showService) CreateShows(request model.Request) error {
	return nil
}

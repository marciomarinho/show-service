package service

import (
	"github.com/marciomarinho/show-service/internal/domain"
	"github.com/marciomarinho/show-service/internal/repository"
)

type ShowService interface {
	Create(show domain.Show) error
	Get(slug string) (*domain.Show, error)
	List() ([]domain.Show, error)
}

type ShowSvc struct {
	repo repository.ShowRepository
}

func NewShowService(r repository.ShowRepository) ShowService {
	return &ShowSvc{repo: r}
}

func (s *ShowSvc) Create(show domain.Show) error {
	return s.repo.Put(show)
}

func (s *ShowSvc) Get(slug string) (*domain.Show, error) {
	return s.repo.Get(slug)
}

func (s *ShowSvc) List() ([]domain.Show, error) {
	return s.repo.List()
}

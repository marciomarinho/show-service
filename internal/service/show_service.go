package service

import "github.com/marciomarinho/show-service/internal/domain"

type ShowService struct {
	repo domain.ShowRepository
}

func NewShowService(r domain.ShowRepository) *ShowService {
	return &ShowService{repo: r}
}

func (s *ShowService) Create(show domain.Show) error {
	return s.repo.Put(show)
}

func (s *ShowService) Get(slug string) (*domain.Show, error) {
	return s.repo.Get(slug)
}

func (s *ShowService) List() ([]domain.Show, error) {
	return s.repo.List()
}

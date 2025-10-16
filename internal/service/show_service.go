package service

import (
	"github.com/marciomarinho/show-service/internal/domain"
	"github.com/marciomarinho/show-service/internal/repository"
)

type ShowService interface {
	Create(request domain.Request) error
	List() (*domain.Response, error)
}

type ShowSvc struct {
	repo repository.ShowRepository
}

func NewShowService(r repository.ShowRepository) ShowService {
	return &ShowSvc{repo: r}
}

func (s *ShowSvc) Create(request domain.Request) error {
	// Save all shows in the request payload
	for _, show := range request.Payload {
		if err := s.repo.Put(show); err != nil {
			return err
		}
	}
	return nil
}

func (s *ShowSvc) List() (*domain.Response, error) {
	shows, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	// Convert domain.Show to domain.ShowResponse for API response
	var showResponses []domain.ShowResponse
	for _, show := range shows {
		showResponse := domain.ShowResponse{
			Image: getImageURL(show.Image),
			Slug:  show.Slug,
			Title: show.Title,
		}
		showResponses = append(showResponses, showResponse)
	}

	response := &domain.Response{
		Response: showResponses,
	}

	return response, nil
}

func getImageURL(img *domain.Image) string {
	if img == nil {
		return ""
	}
	return img.ShowImage
}

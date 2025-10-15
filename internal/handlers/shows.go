package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marciomarinho/show-service/internal/domain"
	"github.com/marciomarinho/show-service/internal/service"
)

type ShowHTTP struct {
	svc *service.ShowService
}

func NewShowHTTP(s *service.ShowService) *ShowHTTP {
	return &ShowHTTP{svc: s}
}

func (h *ShowHTTP) PostShows(c *gin.Context) {
	var req domain.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For now, just handle the first show in the payload
	// In a real implementation, you'd handle batch creation
	if len(req.Payload) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no shows in payload"})
		return
	}

	show := req.Payload[0]
	if err := h.svc.Create(show); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *ShowHTTP) GetShows(c *gin.Context) {
	shows, err := h.svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	response := domain.Response{
		Response: showResponses,
	}

	c.JSON(http.StatusOK, response)
}

// getImageURL extracts image URL from the Image struct
func getImageURL(img *domain.Image) string {
	if img == nil {
		return ""
	}
	return img.ShowImage
}

func (h *ShowHTTP) GetShowBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug parameter is required"})
		return
	}

	show, err := h.svc.Get(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if show == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "show not found"})
		return
	}

	// Convert domain.Show to domain.ShowResponse for API response
	showResponse := domain.ShowResponse{
		Image: getImageURL(show.Image),
		Slug:  show.Slug,
		Title: show.Title,
	}

	c.JSON(http.StatusOK, showResponse)
}

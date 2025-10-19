package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marciomarinho/show-service/internal/domain"
	"github.com/marciomarinho/show-service/internal/service"
)

type ShowHandler interface {
	PostShows(c *gin.Context)
	GetShows(c *gin.Context)
}

type ShowHTTPHandler struct {
	svc service.ShowService
}

func NewShowHandler(s service.ShowService) ShowHandler {
	return &ShowHTTPHandler{svc: s}
}

func (h *ShowHTTPHandler) PostShows(c *gin.Context) {
	var req domain.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not decode request: " + err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not decode request: " + err.Error()})
		return
	}

	for _, show := range req.Payload {
		if err := show.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not decode request: " + err.Error()})
			return
		}
	}

	if err := h.svc.Create(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Shows created successfully"})
}

func (h *ShowHTTPHandler) GetShows(c *gin.Context) {
	response, err := h.svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

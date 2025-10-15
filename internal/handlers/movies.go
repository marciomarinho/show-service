package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	model "github.com/stan-projects/show-service/internal/model"
)

// getAlbums responds with the list of all albums as JSON.
func PostShows(c *gin.Context) {

	var request model.Request
	if err := c.BindJSON(&request); err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "created"})
}

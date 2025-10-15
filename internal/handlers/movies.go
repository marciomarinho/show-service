package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/marciomarinho/show-service/internal/model"
)

// getAlbums responds with the list of all albums as JSON.
func PostShows(c *gin.Context) {

	var request model.Request
	if err := c.BindJSON(&request); err != nil {
		return
	}

	//c.IndentedJSON(http.StatusCreated, gin.H{"message": "created"})
	c.IndentedJSON(http.StatusCreated, request)
}

// getAlbums responds with the list of all albums as JSON.
func GetShows(c *gin.Context) {

	response := &model.Response{
		Response: []model.ShowResponse{
			{
				Image: "http://catchup.ninemsn.com.au/img/jump-in/shows/16KidsandCounting1280.jpg",
				Slug:  "show/16kidsandcounting",
				Title: "16 Kids and Counting",
			},
			{
				Image: "http://catchup.ninemsn.com.au/img/jump-in/shows/TheTaste1280.jpg",
				Slug:  "show/thetaste",
				Title: "The Taste (Le Goût)",
			},
			{
				Image: "http://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg",
				Slug:  "show/thunderbirds",
				Title: "Thunderbirds",
			},
			{
				Image: "http://catchup.ninemsn.com.au/img/jump-in/shows/ScoobyDoo1280.jpg",
				Slug:  "show/scoobydoomysteryincorporated",
				Title: "Scooby-Doo! Mystery Incorporated",
			},
			{
				Image: "http://catchup.ninemsn.com.au/img/jump-in/shows/ToyHunter1280.jpg",
				Slug:  "show/toyhunter",
				Title: "Toy Hunter",
			},
			{
				Image: "http://catchup.ninemsn.com.au/img/jump-in/shows/Worlds1280.jpg",
				Slug:  "show/worlds",
				Title: "World's...",
			},
			{
				Image: "http://catchup.ninemsn.com.au/img/jump-in/shows/TheOriginals1280.jpg",
				Slug:  "show/theoriginals",
				Title: "The Originals",
			},
		},
	}

	//c.IndentedJSON(http.StatusCreated, gin.H{"message": "created"})
	c.IndentedJSON(http.StatusCreated, response)
}

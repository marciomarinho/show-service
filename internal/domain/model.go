package domain

// Image represents show image information
type Image struct {
	ShowImage string `json:"showImage" dynamodbav:"showImage"`
}

// NextEpisode represents next episode information
type NextEpisode struct {
	Channel     *string `json:"channel,omitempty" dynamodbav:"channel"`
	ChannelLogo string  `json:"channelLogo" dynamodbav:"channelLogo"`
	Date        *string `json:"date,omitempty" dynamodbav:"date"`
	HTML        string  `json:"html" dynamodbav:"html"`
	URL         string  `json:"url" dynamodbav:"url"`
}

// Season represents a show season
type Season struct {
	Slug string `json:"slug" dynamodbav:"slug"`
}

// Show represents a complete show model
type Show struct {
	Country       *string      `json:"country,omitempty" dynamodbav:"country"`
	Description   *string      `json:"description,omitempty" dynamodbav:"description"`
	DRM           *bool        `json:"drm,omitempty" dynamodbav:"drm"`
	EpisodeCount  *int         `json:"episodeCount,omitempty" dynamodbav:"episodeCount"`
	Genre         *string      `json:"genre,omitempty" dynamodbav:"genre"`
	Image         *Image       `json:"image,omitempty" dynamodbav:"image"`
	Language      *string      `json:"language,omitempty" dynamodbav:"language"`
	NextEpisode   *NextEpisode `json:"nextEpisode,omitempty" dynamodbav:"nextEpisode"`
	PrimaryColour *string      `json:"primaryColour,omitempty" dynamodbav:"primaryColour"`
	Seasons       *[]Season    `json:"seasons,omitempty" dynamodbav:"seasons"`
	Slug          string       `json:"slug" dynamodbav:"slug"` // PK
	Title         string       `json:"title" dynamodbav:"title"`
	TVChannel     *string      `json:"tvChannel,omitempty" dynamodbav:"tvChannel"`
}

// Request represents API request with pagination
type Request struct {
	Payload      []Show `json:"payload"`
	Skip         int    `json:"skip"`
	Take         int    `json:"take"`
	TotalRecords int    `json:"totalRecords"`
}

// Response represents API response
type Response struct {
	Response []ShowResponse `json:"response"`
}

// ShowResponse represents a simplified show for API responses
type ShowResponse struct {
	Image string `json:"image" dynamodbav:"image"`
	Slug  string `json:"slug" dynamodbav:"slug"`
	Title string `json:"title" dynamodbav:"title"`
}

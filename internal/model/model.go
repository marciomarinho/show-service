package model

//Response model

type Request struct {
	Payload      []Show `json:"payload"`
	Skip         int    `json:"skip"`
	Take         int    `json:"take"`
	TotalRecords int    `json:"totalRecords"`
}

type Show struct {
	Country       *string      `json:"country,omitempty"`
	Description   *string      `json:"description,omitempty"`
	DRM           *bool        `json:"drm,omitempty"`
	EpisodeCount  *int         `json:"episodeCount,omitempty"`
	Genre         *string      `json:"genre,omitempty"`
	Image         *Image       `json:"image,omitempty"`
	Language      *string      `json:"language,omitempty"`
	NextEpisode   *NextEpisode `json:"nextEpisode,omitempty"`
	PrimaryColour *string      `json:"primaryColour,omitempty"`
	Seasons       *[]Season    `json:"seasons,omitempty"`
	Slug          string       `json:"slug"`
	Title         string       `json:"title"`
	TVChannel     *string      `json:"tvChannel,omitempty"`
}

type Image struct {
	ShowImage string `json:"showImage"`
}

type NextEpisode struct {
	Channel     *string `json:"channel,omitempty"`
	ChannelLogo string  `json:"channelLogo"`
	Date        *string `json:"date,omitempty"`
	HTML        string  `json:"html"`
	URL         string  `json:"url"`
}

type Season struct {
	Slug string `json:"slug"`
}

//Response model

type Response struct {
	Response []ShowResponse `json:"response"`
}

type ShowResponse struct {
	Image string `json:"image"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

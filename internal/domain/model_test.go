package domain

import (
	"testing"
)

func TestImage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		image   Image
		wantErr bool
	}{
		{
			name: "valid image",
			image: Image{
				ShowImage: "http://example.com/image.jpg",
			},
			wantErr: false,
		},
		{
			name: "valid image with query params",
			image: Image{
				ShowImage: "https://cdn.example.com/img/show.jpg?w=1280&h=720",
			},
			wantErr: false,
		},
		{
			name: "valid image with port",
			image: Image{
				ShowImage: "http://localhost:8080/image.png",
			},
			wantErr: false,
		},
		{
			name: "valid image with subdomain",
			image: Image{
				ShowImage: "https://images.example.com/show/16KidsandCounting1280.jpg",
			},
			wantErr: false,
		},
		{
			name: "missing showImage",
			image: Image{
				ShowImage: "",
			},
			wantErr: true,
		},
		{
			name: "invalid URL format",
			image: Image{
				ShowImage: "not-a-url",
			},
			wantErr: true,
		},
		{
			name: "invalid URL scheme",
			image: Image{
				ShowImage: "ftp://example.com/image.jpg",
			},
			wantErr: true,
		},
		{
			name: "URL with invalid characters",
			image: Image{
				ShowImage: "http://example.com/image with spaces.jpg",
			},
			wantErr: true,
		},
		{
			name: "malformed URL",
			image: Image{
				ShowImage: "http:///missing-host",
			},
			wantErr: true,
		},
		{
			name: "URL with special characters",
			image: Image{
				ShowImage: "http://example.com/image%20with%20spaces.jpg",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.image.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Image.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNextEpisode_Validate(t *testing.T) {
	tests := []struct {
		name        string
		nextEpisode NextEpisode
		wantErr     bool
	}{
		{
			name: "valid next episode",
			nextEpisode: NextEpisode{
				ChannelLogo: "http://example.com/logo.gif",
				HTML:        "Next episode airs: 10:00pm Monday",
				URL:         "http://example.com/episode",
			},
			wantErr: false,
		},
		{
			name: "valid next episode with optional fields",
			nextEpisode: NextEpisode{
				Channel:     stringPtr("Channel 9"),
				ChannelLogo: "http://example.com/logo.gif",
				Date:        stringPtr("2023-10-16"),
				HTML:        "Next episode airs: 10:00pm Monday",
				URL:         "http://example.com/episode",
			},
			wantErr: false,
		},
		{
			name: "missing channelLogo",
			nextEpisode: NextEpisode{
				HTML: "Next episode airs: 10:00pm Monday",
				URL:  "http://example.com/episode",
			},
			wantErr: true,
		},
		{
			name: "missing HTML",
			nextEpisode: NextEpisode{
				ChannelLogo: "http://example.com/logo.gif",
				URL:         "http://example.com/episode",
			},
			wantErr: true,
		},
		{
			name: "missing URL",
			nextEpisode: NextEpisode{
				ChannelLogo: "http://example.com/logo.gif",
				HTML:        "Next episode airs: 10:00pm Monday",
			},
			wantErr: true,
		},
		{
			name: "invalid URL format",
			nextEpisode: NextEpisode{
				ChannelLogo: "http://example.com/logo.gif",
				HTML:        "Next episode airs: 10:00pm Monday",
				URL:         "not-a-valid-url",
			},
			wantErr: true,
		},
		{
			name: "empty channelLogo URL",
			nextEpisode: NextEpisode{
				ChannelLogo: "",
				HTML:        "Next episode airs: 10:00pm Monday",
				URL:         "http://example.com/episode",
			},
			wantErr: true,
		},
		{
			name: "empty HTML",
			nextEpisode: NextEpisode{
				ChannelLogo: "http://example.com/logo.gif",
				HTML:        "",
				URL:         "http://example.com/episode",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.nextEpisode.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("NextEpisode.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSeason_Validate(t *testing.T) {
	tests := []struct {
		name    string
		season  Season
		wantErr bool
	}{
		{
			name: "valid season",
			season: Season{
				Slug: "show/16kidsandcounting/season/1",
			},
			wantErr: false,
		},
		{
			name: "valid season with multiple words",
			season: Season{
				Slug: "show/scoobydoomysteryincorporated/season/1",
			},
			wantErr: false,
		},
		{
			name: "valid season with numbers",
			season: Season{
				Slug: "show/thunderbirds/season/8",
			},
			wantErr: false,
		},
		{
			name: "missing slug",
			season: Season{
				Slug: "",
			},
			wantErr: true,
		},
		{
			name: "invalid slug format - missing show prefix",
			season: Season{
				Slug: "16kidsandcounting/season/1",
			},
			wantErr: true,
		},
		{
			name: "invalid slug format - uppercase",
			season: Season{
				Slug: "SHOW/16kidsandcounting/season/1",
			},
			wantErr: true,
		},
		{
			name: "invalid slug format - spaces",
			season: Season{
				Slug: "show/16 kids and counting/season/1",
			},
			wantErr: true,
		},
		{
			name: "invalid slug format - special characters",
			season: Season{
				Slug: "show/16kidsandcounting/season/1!",
			},
			wantErr: true,
		},
		{
			name: "valid minimal slug",
			season: Season{
				Slug: "show/test",
			},
			wantErr: false,
		},
		{
			name: "valid slug with underscore",
			season: Season{
				Slug: "show/test_show",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.season.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Season.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShow_Validate(t *testing.T) {
	tests := []struct {
		name    string
		show    Show
		wantErr bool
	}{
		{
			name: "valid minimal show",
			show: Show{
				Slug:    "show/testshow",
				Title:   "Test Show",
				Seasons: &[]Season{}, // Initialize as empty slice
			},
			wantErr: false,
		},
		{
			name: "valid complete show",
			show: Show{
				Country:      stringPtr("USA"),
				Description:  stringPtr("A test show description"),
				DRM:          boolPtr(true),
				EpisodeCount: intPtr(10),
				Genre:        stringPtr("Comedy"),
				Image: &Image{
					ShowImage: "http://example.com/image.jpg",
				},
				Language:      stringPtr("English"),
				PrimaryColour: stringPtr("#ff7800"),
				Seasons: &[]Season{
					{Slug: "show/testshow/season/1"},
				},
				Slug:      "show/testshow",
				Title:     "Test Show",
				TVChannel: stringPtr("Channel 9"),
			},
			wantErr: false,
		},
		{
			name: "show with next episode",
			show: Show{
				Slug:  "show/testshow",
				Title: "Test Show",
				NextEpisode: &NextEpisode{
					ChannelLogo: "http://example.com/logo.gif",
					HTML:        "Next episode airs tomorrow",
					URL:         "http://example.com/episode",
				},
				Seasons: &[]Season{},
			},
			wantErr: false,
		},
		{
			name: "missing slug",
			show: Show{
				Title:   "Test Show",
				Seasons: &[]Season{},
			},
			wantErr: true,
		},
		{
			name: "missing title",
			show: Show{
				Slug:    "show/testshow",
				Seasons: &[]Season{},
			},
			wantErr: true,
		},
		{
			name: "invalid slug format",
			show: Show{
				Slug:    "invalid-slug",
				Title:   "Test Show",
				Seasons: &[]Season{},
			},
			wantErr: true,
		},
		{
			name: "title too long",
			show: Show{
				Slug:    "show/testshow",
				Title:   "This is a very long title that exceeds the maximum allowed length of two hundred characters and should fail validation because it's way too long for the database and API response constraints",
				Seasons: &[]Season{},
			},
			wantErr: true,
		},
		{
			name: "invalid hex color",
			show: Show{
				Slug:          "show/testshow",
				Title:         "Test Show",
				PrimaryColour: stringPtr("red"),
				Seasons:       &[]Season{},
			},
			wantErr: true,
		},
		{
			name: "invalid hex color format",
			show: Show{
				Slug:          "show/testshow",
				Title:         "Test Show",
				PrimaryColour: stringPtr("#12345"),
				Seasons:       &[]Season{},
			},
			wantErr: true,
		},
		{
			name: "negative episode count",
			show: Show{
				Slug:         "show/testshow",
				Title:        "Test Show",
				EpisodeCount: intPtr(-1),
				Seasons:      &[]Season{},
			},
			wantErr: true,
		},
		{
			name: "country too long",
			show: Show{
				Slug:    "show/testshow",
				Title:   "Test Show",
				Country: stringPtr("This is a very long country name that exceeds the maximum allowed length"),
				Seasons: &[]Season{},
			},
			wantErr: true,
		},
		{
			name: "description too long",
			show: Show{
				Slug:        "show/testshow",
				Title:       "Test Show",
				Description: stringPtr("This is a very long description that exceeds the maximum allowed length of 500 characters and should fail validation because it's way too long for storage and display purposes in the application interface. This description is being extended with additional text to ensure it surpasses the 500 character limit for testing the validation logic in the Show model's Validate method. The test case 'description too long' requires a string longer than 500 characters to verify that the validation correctly identifies and rejects overly long descriptions."),
				Seasons:     &[]Season{},
			},
			wantErr: true,
		},
		{
			name: "genre too long",
			show: Show{
				Slug:    "show/testshow",
				Title:   "Test Show",
				Genre:   stringPtr("This is a very long genre name that exceeds the limit"),
				Seasons: &[]Season{},
			},
			wantErr: true,
		},
		{
			name: "language too long",
			show: Show{
				Slug:     "show/testshow",
				Title:    "Test Show",
				Language: stringPtr("This is a very long language name that exceeds the limit"),
				Seasons:  &[]Season{},
			},
			wantErr: true,
		},
		{
			name: "tvChannel too long",
			show: Show{
				Slug:      "show/testshow",
				Title:     "Test Show",
				TVChannel: stringPtr("This is a very long TV channel name that exceeds the maximum allowed length"),
				Seasons:   &[]Season{},
			},
			wantErr: true,
		},
		{
			name: "invalid image URL",
			show: Show{
				Slug:    "show/testshow",
				Title:   "Test Show",
				Image:   &Image{ShowImage: "not-a-valid-url"},
				Seasons: &[]Season{},
			},
			wantErr: true,
		},
		{
			name: "invalid next episode",
			show: Show{
				Slug:  "show/testshow",
				Title: "Test Show",
				NextEpisode: &NextEpisode{
					ChannelLogo: "http://example.com/logo.gif",
					HTML:        "Next episode",
					URL:         "invalid-url",
				},
				Seasons: &[]Season{},
			},
			wantErr: true,
		},
		{
			name: "empty seasons array",
			show: Show{
				Slug:    "show/testshow",
				Title:   "Test Show",
				Seasons: &[]Season{},
			},
			wantErr: false,
		},
		{
			name: "seasons with invalid slug",
			show: Show{
				Slug:  "show/testshow",
				Title: "Test Show",
				Seasons: &[]Season{
					{Slug: "invalid-season-slug"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.show.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Show.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request Request
		wantErr bool
	}{
		{
			name: "valid request",
			request: Request{
				Payload: []Show{
					{Slug: "show/testshow1", Title: "Test Show 1", Seasons: &[]Season{}},
					{Slug: "show/testshow2", Title: "Test Show 2", Seasons: &[]Season{}},
				},
				Skip:         0,
				Take:         10,
				TotalRecords: 50,
			},
			wantErr: false,
		},
		{
			name: "valid request with maximum payload",
			request: Request{
				Payload: func() []Show {
					shows := make([]Show, 1000)
					for i := range shows {
						shows[i] = Show{
							Slug:    "show/testshow",
							Title:   "Test Show",
							Seasons: &[]Season{},
						}
					}
					return shows
				}(),
				Skip:         0,
				Take:         100,
				TotalRecords: 1000,
			},
			wantErr: false,
		},
		{
			name: "empty payload",
			request: Request{
				Payload:      []Show{},
				Skip:         0,
				Take:         10,
				TotalRecords: 0,
			},
			wantErr: true,
		},
		{
			name: "payload too large",
			request: Request{
				Payload: func() []Show {
					shows := make([]Show, 1001)
					for i := range shows {
						shows[i] = Show{
							Slug:    "show/testshow",
							Title:   "Test Show",
							Seasons: &[]Season{},
						}
					}
					return shows
				}(),
				Skip:         0,
				Take:         10,
				TotalRecords: 1001,
			},
			wantErr: true,
		},
		{
			name: "negative skip",
			request: Request{
				Payload:      []Show{{Slug: "show/test", Title: "Test", Seasons: &[]Season{}}},
				Skip:         -1,
				Take:         10,
				TotalRecords: 1,
			},
			wantErr: true,
		},
		{
			name: "zero take",
			request: Request{
				Payload:      []Show{{Slug: "show/test", Title: "Test", Seasons: &[]Season{}}},
				Skip:         0,
				Take:         0,
				TotalRecords: 1,
			},
			wantErr: true,
		},
		{
			name: "take too large",
			request: Request{
				Payload:      []Show{{Slug: "show/test", Title: "Test", Seasons: &[]Season{}}},
				Skip:         0,
				Take:         101, // Exceeds max
				TotalRecords: 1,
			},
			wantErr: true,
		},
		{
			name: "negative total records",
			request: Request{
				Payload:      []Show{{Slug: "show/test", Title: "Test", Seasons: &[]Season{}}},
				Skip:         0,
				Take:         10,
				TotalRecords: -1,
			},
			wantErr: true,
		},
		{
			name: "single item payload",
			request: Request{
				Payload:      []Show{{Slug: "show/test", Title: "Test", Seasons: &[]Season{}}},
				Skip:         0,
				Take:         1,
				TotalRecords: 1,
			},
			wantErr: false,
		},
		{
			name: "request with invalid show",
			request: Request{
				Payload: []Show{
					{Slug: "invalid-slug", Title: "Test", Seasons: &[]Season{}}, // Invalid slug
				},
				Skip:         0,
				Take:         10,
				TotalRecords: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Request.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper functions for tests
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

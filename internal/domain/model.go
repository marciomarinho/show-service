package domain

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// MatchShowSlug validates show slug format: show/<handle>
// handle: letters/digits/dashes, must start with letter or digit
var MatchShowSlug = regexp.MustCompile(`^show/[a-z0-9][a-z0-9-]*$`)

// MatchSeasonSlug validates season slug format: show/<handle>/season/<number>=1+
var MatchSeasonSlug = regexp.MustCompile(`^show/[a-z0-9][a-z0-9-]*/season/[1-9][0-9]*$`)

// MatchHexColor validates hex color format (#ffffff)
var MatchHexColor = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

// ValidateURL validates URL format
func ValidateURL(value interface{}) error {
	s, _ := value.(string)
	if s == "" {
		return nil
	}

	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		return validation.NewError("invalid_url", "must start with http:// or https://")
	}

	parsed, err := url.Parse(s)
	if err != nil {
		return validation.NewError("invalid_url", "must be a valid URL")
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return validation.NewError("invalid_url", "must have valid scheme and host")
	}

	if strings.Contains(s, " ") {
		return validation.NewError("invalid_url", "URL cannot contain spaces")
	}

	return nil
}

type URLRule struct{}

func (r URLRule) Validate(value interface{}) error {
	return ValidateURL(value)
}

func ValidateStringLength(s *string, min, max int) error {
	if s == nil {
		return nil
	}
	if len(*s) < min || len(*s) > max {
		return validation.NewError("length_error", "string length out of range")
	}
	return nil
}

type Image struct {
	ShowImage string `json:"showImage" dynamodbav:"showImage"`
}

func (i Image) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.ShowImage, validation.Required, URLRule{}),
	)
}

// NextEpisode represents next episode information
type NextEpisode struct {
	Channel     *string `json:"channel,omitempty" dynamodbav:"channel"`
	ChannelLogo string  `json:"channelLogo" dynamodbav:"channelLogo"`
	Date        *string `json:"date,omitempty" dynamodbav:"date"`
	HTML        string  `json:"html" dynamodbav:"html"`
	URL         string  `json:"url" dynamodbav:"url"`
}

func (n NextEpisode) Validate() error {
	return validation.ValidateStruct(&n,
		validation.Field(&n.ChannelLogo, validation.Required),
		validation.Field(&n.HTML, validation.Required),
		validation.Field(&n.URL, validation.Required, URLRule{}),
	)
}

type Season struct {
	Slug string `json:"slug" dynamodbav:"slug"`
}

func (s Season) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Slug, validation.Required, validation.By(func(value interface{}) error {
			v, _ := value.(string)
			if MatchSeasonSlug.MatchString(v) || MatchShowSlug.MatchString(v) {
				return nil
			}
			return validation.NewError("invalid_slug", "must be in a valid format")
		})),
	)
}

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

	// Index helpers (not in JSON payloads; set on write for GSI)
	DRMKey *int `json:"-" dynamodbav:"drmKey,omitempty"`
}

func (s Show) Validate() error {
	if len(strings.TrimSpace(s.Title)) == 0 {
		return validation.NewError("title_required", "title is required")
	}
	if len(s.Title) > 120 {
		return validation.NewError("title_too_long", "title must be at most 120 characters")
	}

	if err := ValidateStringLength(s.Country, 0, 50); err != nil {
		return fmt.Errorf("country: %w", err)
	}
	if err := ValidateStringLength(s.Genre, 0, 50); err != nil {
		return fmt.Errorf("genre: %w", err)
	}
	if err := ValidateStringLength(s.Language, 0, 50); err != nil {
		return fmt.Errorf("language: %w", err)
	}
	if err := ValidateStringLength(s.TVChannel, 0, 50); err != nil {
		return fmt.Errorf("tvChannel: %w", err)
	}
	if err := ValidateStringLength(s.Description, 0, 500); err != nil {
		return fmt.Errorf("description: %w", err)
	}

	return validation.ValidateStruct(&s,
		validation.Field(&s.Slug, validation.Required, validation.Match(MatchShowSlug)),
		validation.Field(&s.PrimaryColour, validation.When(s.PrimaryColour != nil, validation.Match(MatchHexColor).Error("must be valid hex color"))),
		validation.Field(&s.EpisodeCount, validation.When(s.EpisodeCount != nil, validation.Min(0))),
		validation.Field(&s.Image),
		validation.Field(&s.NextEpisode),
		validation.Field(&s.Seasons, validation.When(s.Seasons != nil,
			validation.By(func(value any) error {
				seasons := value.(*[]Season)
				if seasons == nil {
					return nil
				}
				return validation.Validate(*seasons, validation.Each())
			}),
		)),
	)
}

type Request struct {
	Payload      []Show `json:"payload"`
	Skip         int    `json:"skip"`
	Take         int    `json:"take"`
	TotalRecords int    `json:"totalRecords"`
}

func (r Request) Validate() error {
	if len(r.Payload) < 1 || len(r.Payload) > 1000 {
		return validation.NewError("payload_size", "payload must contain between 1 and 1000 items")
	}
	if r.Skip < 0 {
		return validation.NewError("skip_invalid", "skip must be >= 0")
	}
	if r.Take < 1 || r.Take > 100 {
		return validation.NewError("take_invalid", "take must be between 1 and 100")
	}
	if r.TotalRecords < 0 {
		return validation.NewError("total_records_invalid", "totalRecords must be >= 0")
	}

	for i := range r.Payload {
		if err := r.Payload[i].Validate(); err != nil {
			return fmt.Errorf("payload[%d]: %w", i, err)
		}
	}
	return nil
}

type Response struct {
	Response []ShowResponse `json:"response"`
}

type ShowResponse struct {
	Image string `json:"image" dynamodbav:"image"`
	Slug  string `json:"slug" dynamodbav:"slug"`
	Title string `json:"title" dynamodbav:"title"`
}

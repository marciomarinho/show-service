package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/marciomarinho/show-service/internal/domain"
	"github.com/marciomarinho/show-service/internal/repository"
)

func TestShowSvc_Create(t *testing.T) {
	tests := []struct {
		name        string
		request     domain.Request
		mockSetup   func(*repository.MockShowRepository)
		expectError bool
	}{
		{
			name: "successful creation of multiple shows",
			request: domain.Request{
				Payload: []domain.Show{
					{
						Slug:  "show/test1",
						Title: "Test Show 1",
						DRM:   &[]bool{true}[0],
					},
					{
						Slug:  "show/test2",
						Title: "Test Show 2",
						DRM:   &[]bool{false}[0],
					},
				},
			},
			mockSetup: func(m *repository.MockShowRepository) {
				m.On("Put", mock.AnythingOfType("domain.Show")).Return(nil).Times(2)
			},
			expectError: false,
		},
		{
			name: "creation fails on first show",
			request: domain.Request{
				Payload: []domain.Show{
					{
						Slug:  "show/test1",
						Title: "Test Show 1",
						DRM:   &[]bool{true}[0],
					},
					{
						Slug:  "show/test2",
						Title: "Test Show 2",
						DRM:   &[]bool{false}[0],
					},
				},
			},
			mockSetup: func(m *repository.MockShowRepository) {
				m.On("Put", mock.AnythingOfType("domain.Show")).Return(errors.New("database error")).Once()
			},
			expectError: true,
		},
		{
			name: "creation fails on second show",
			request: domain.Request{
				Payload: []domain.Show{
					{
						Slug:  "show/test1",
						Title: "Test Show 1",
						DRM:   &[]bool{true}[0],
					},
					{
						Slug:  "show/test2",
						Title: "Test Show 2",
						DRM:   &[]bool{false}[0],
					},
				},
			},
			mockSetup: func(m *repository.MockShowRepository) {
				m.On("Put", mock.AnythingOfType("domain.Show")).Return(nil).Once()
				m.On("Put", mock.AnythingOfType("domain.Show")).Return(errors.New("database error")).Once()
			},
			expectError: true,
		},
		{
			name: "empty payload",
			request: domain.Request{
				Payload: []domain.Show{},
			},
			mockSetup: func(m *repository.MockShowRepository) {
				// No calls expected for empty payload
			},
			expectError: false,
		},
		{
			name: "single show creation",
			request: domain.Request{
				Payload: []domain.Show{
					{
						Slug:  "show/single",
						Title: "Single Show",
						DRM:   &[]bool{true}[0],
					},
				},
			},
			mockSetup: func(m *repository.MockShowRepository) {
				m.On("Put", mock.AnythingOfType("domain.Show")).Return(nil).Once()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockShowRepository)
			tt.mockSetup(mockRepo)

			svc := NewShowService(mockRepo)
			err := svc.Create(tt.request)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestShowSvc_List(t *testing.T) {
	tests := []struct {
		name        string
		mockShows   []domain.Show
		mockError   error
		expectError bool
		expectedLen int
	}{
		{
			name: "successful list with shows",
			mockShows: []domain.Show{
				{
					Slug:  "show/test1",
					Title: "Test Show 1",
					Image: &domain.Image{ShowImage: "http://example.com/image1.jpg"},
				},
				{
					Slug:  "show/test2",
					Title: "Test Show 2",
					Image: &domain.Image{ShowImage: "http://example.com/image2.jpg"},
				},
			},
			mockError:   nil,
			expectError: false,
			expectedLen: 2,
		},
		{
			name:        "successful list with no shows",
			mockShows:   []domain.Show{},
			mockError:   nil,
			expectError: false,
			expectedLen: 0,
		},
		{
			name: "repository error",
			mockShows: []domain.Show{
				{
					Slug:  "show/test1",
					Title: "Test Show 1",
				},
			},
			mockError:   errors.New("database connection error"),
			expectError: true,
			expectedLen: 0,
		},
		{
			name: "single show",
			mockShows: []domain.Show{
				{
					Slug:  "show/single",
					Title: "Single Show",
					Image: nil,
				},
			},
			mockError:   nil,
			expectError: false,
			expectedLen: 1,
		},
		{
			name: "show with nil image",
			mockShows: []domain.Show{
				{
					Slug:  "show/noimage",
					Title: "Show Without Image",
					Image: nil,
				},
			},
			mockError:   nil,
			expectError: false,
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockShowRepository)
			mockRepo.On("List").Return(tt.mockShows, tt.mockError)

			svc := NewShowService(mockRepo)
			response, err := svc.List()

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, response)
			} else {
				require.NoError(t, err)
				require.NotNil(t, response)
				require.Len(t, response.Response, tt.expectedLen)

				for i, showResp := range response.Response {
					require.Equal(t, tt.mockShows[i].Slug, showResp.Slug)
					require.Equal(t, tt.mockShows[i].Title, showResp.Title)

					if tt.mockShows[i].Image != nil {
						require.Equal(t, tt.mockShows[i].Image.ShowImage, showResp.Image)
					} else {
						require.Equal(t, "", showResp.Image)
					}
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

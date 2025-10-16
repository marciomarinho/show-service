package repository

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/marciomarinho/show-service/internal/domain"
	"github.com/marciomarinho/show-service/internal/repository/mocks"
)

func TestShowRepo_Put(t *testing.T) {
	t.Run("valid show insertion", func(t *testing.T) {
		mockDB := mocks.NewMockDynamoAPI(t)

		// Set up expectations
		mockDB.On("TableName").Return("test-table").Maybe()
		mockDB.On("PutItem", mock.Anything, mock.AnythingOfType("*dynamodb.PutItemInput")).
			Run(func(args mock.Arguments) {
				in := args.Get(1).(*dynamodb.PutItemInput)
				require.Equal(t, "test-table", *in.TableName)
				require.NotEmpty(t, in.Item["slug"])
			}).
			Return(&dynamodb.PutItemOutput{}, nil)

		repo := NewShowRepository(mockDB)

		err := repo.Put(domain.Show{
			Slug:    "show/testshow",
			Title:   "Test Show",
			DRM:     boolPtr(true),
			Seasons: &[]domain.Season{},
		})
		require.NoError(t, err)
	})

	t.Run("empty slug error", func(t *testing.T) {
		mockDB := mocks.NewMockDynamoAPI(t)
		repo := NewShowRepository(mockDB)

		err := repo.Put(domain.Show{
			Slug:    "", // invalid
			Title:   "Test Show",
			Seasons: &[]domain.Season{},
		})
		require.Error(t, err)
		// PutItem should never be called
		mockDB.AssertNotCalled(t, "PutItem", mock.Anything, mock.Anything)
	})

	t.Run("validation error", func(t *testing.T) {
		mockDB := mocks.NewMockDynamoAPI(t)
		repo := NewShowRepository(mockDB)

		err := repo.Put(domain.Show{
			Slug:    "show/testshow",
			Title:   "", // invalid - empty title
			Seasons: &[]domain.Season{},
		})
		require.Error(t, err)
		// PutItem should never be called due to validation
		mockDB.AssertNotCalled(t, "PutItem", mock.Anything, mock.Anything)
	})
}

func TestShowRepo_List(t *testing.T) {
	t.Run("successful list with DRM shows", func(t *testing.T) {
		mockDB := mocks.NewMockDynamoAPI(t)

		mockDB.On("TableName").Return("test-table").Maybe()
		mockDB.On("Scan", mock.Anything, mock.AnythingOfType("*dynamodb.ScanInput")).
			Return(&dynamodb.ScanOutput{
				Items: []map[string]types.AttributeValue{
					{
						"slug":         &types.AttributeValueMemberS{Value: "show/a"},
						"title":        &types.AttributeValueMemberS{Value: "A"},
						"drm":          &types.AttributeValueMemberBOOL{Value: true},
						"episodeCount": &types.AttributeValueMemberN{Value: "3"},
						"image": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
							"showImage": &types.AttributeValueMemberS{Value: "http://x/y.jpg"},
						}},
					},
					{
						"slug":         &types.AttributeValueMemberS{Value: "show/b"},
						"title":        &types.AttributeValueMemberS{Value: "B"},
						"drm":          &types.AttributeValueMemberBOOL{Value: true},
						"episodeCount": &types.AttributeValueMemberN{Value: "0"},
					},
				},
			}, nil)

		repo := NewShowRepository(mockDB)

		got, err := repo.List()
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, "show/a", got[0].Slug)
	})

	t.Run("empty result set", func(t *testing.T) {
		mockDB := mocks.NewMockDynamoAPI(t)

		mockDB.On("TableName").Return("test-table").Maybe()
		mockDB.On("Scan", mock.Anything, mock.AnythingOfType("*dynamodb.ScanInput")).
			Return(&dynamodb.ScanOutput{Items: []map[string]types.AttributeValue{}}, nil)

		repo := NewShowRepository(mockDB)

		got, err := repo.List()
		require.NoError(t, err)
		require.Len(t, got, 0)
	})

	t.Run("scan error", func(t *testing.T) {
		mockDB := mocks.NewMockDynamoAPI(t)

		mockDB.On("TableName").Return("test-table").Maybe()
		mockDB.On("Scan", mock.Anything, mock.AnythingOfType("*dynamodb.ScanInput")).
			Return(nil, errors.New("DynamoDB scan failed"))

		repo := NewShowRepository(mockDB)

		got, err := repo.List()
		require.Error(t, err)
		require.Nil(t, got)
	})
}

// helpers
func boolPtr(b bool) *bool { return &b }

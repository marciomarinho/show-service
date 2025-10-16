package repository

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/marciomarinho/show-service/internal/domain"
)

type ShowRepo struct {
	db DynamoAPI
}

func NewShowRepository(db DynamoAPI) *ShowRepo {
	return &ShowRepo{db: db}
}

func (r *ShowRepo) Put(s domain.Show) error {
	// Validate before calling Dynamo
	if err := s.Validate(); err != nil {
		return err
	}

	// derive GSI helpers if you use them (safe no-op if not)
	var k int
	if s.DRM != nil && *s.DRM {
		k = 1
	}
	s.DRMKey = &k
	if s.EpisodeCount == nil {
		zero := 0
		s.EpisodeCount = &zero
	}

	item, err := attributevalue.MarshalMap(s)
	if err != nil {
		return err
	}
	_, err = r.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName:           awsString(r.db.TableName()),
		Item:                item,
		ConditionExpression: awsString("attribute_not_exists(slug)"),
	})
	return err
}

func (r *ShowRepo) List() ([]domain.Show, error) {
	// Query using GSI for DRM=true shows with episodeCount > 0
	// GSI: gsi_drm_episode with hash_key=drmKey, range_key=episodeCount
	out, err := r.db.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              awsString(r.db.TableName()),
		IndexName:              awsString("gsi_drm_episode"),
		KeyConditionExpression: awsString("drmKey = :drmKey AND episodeCount > :episodeCount"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":drmKey":       &types.AttributeValueMemberN{Value: "1"},
			":episodeCount": &types.AttributeValueMemberN{Value: "0"},
		},
	})
	if err != nil {
		return nil, err
	}
	var items []domain.Show
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func awsString(s string) *string { return &s }

// Optional helper to surface clearer errors in tests (not required)
var ErrInvalidShow = errors.New("invalid show")

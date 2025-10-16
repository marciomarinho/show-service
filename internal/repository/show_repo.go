package repository

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
	// simple scan version (your original). Swap to Query if you add a GSI.
	out, err := r.db.Scan(context.Background(), &dynamodb.ScanInput{
		TableName: awsString(r.db.TableName()),
	})
	if err != nil {
		return nil, err
	}
	var items []domain.Show
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &items); err != nil {
		return nil, err
	}
	// filter: drm==true && episodeCount>0
	var filtered []domain.Show
	for _, it := range items {
		if it.DRM != nil && *it.DRM && it.EpisodeCount != nil && *it.EpisodeCount > 0 {
			filtered = append(filtered, it)
		}
	}
	return filtered, nil
}

func awsString(s string) *string { return &s }

// Optional helper to surface clearer errors in tests (not required)
var ErrInvalidShow = errors.New("invalid show")

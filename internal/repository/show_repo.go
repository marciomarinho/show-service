package repository

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/marciomarinho/show-service/internal/domain"
	"github.com/marciomarinho/show-service/internal/infra"
)

type ShowRepository interface {
	Put(show domain.Show) error
	List() ([]domain.Show, error)
}

type ShowRepo struct {
	db    *dynamodb.Client
	table string
}

func NewShowRepository(d *infra.Dynamo) ShowRepository {
	return &ShowRepo{db: d.Client, table: d.TableName}
}

func (r *ShowRepo) Put(show domain.Show) error {
	if show.Slug == "" {
		return errors.New("slug is required")
	}
	item, err := attributevalue.MarshalMap(show)
	if err != nil {
		return err
	}
	_, err = r.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: &r.table, Item: item,
	})
	return err
}

func (r *ShowRepo) List() ([]domain.Show, error) {
	out, err := r.db.Scan(context.Background(), &dynamodb.ScanInput{
		TableName: &r.table,
	})
	if err != nil {
		return nil, err
	}
	var items []domain.Show
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &items); err != nil {
		return nil, err
	}

	// Filter shows: only return those with DRM enabled and at least one episode
	var filteredItems []domain.Show
	for _, item := range items {
		// Check if DRM is enabled and episodeCount > 0
		if item.DRM != nil && *item.DRM && item.EpisodeCount != nil && *item.EpisodeCount > 0 {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems, nil
}

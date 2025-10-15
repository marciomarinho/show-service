package repository

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/marciomarinho/show-service/internal/domain"
	"github.com/marciomarinho/show-service/internal/infra"
)

type ShowRepository interface {
	Put(show domain.Show) error
	Get(slug string) (*domain.Show, error)
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

func (r *ShowRepo) Get(slug string) (*domain.Show, error) {
	out, err := r.db.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"slug": &types.AttributeValueMemberS{Value: slug},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, errors.New("not found")
	}
	var s domain.Show
	if err := attributevalue.UnmarshalMap(out.Item, &s); err != nil {
		return nil, err
	}
	return &s, nil
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
	return items, nil
}

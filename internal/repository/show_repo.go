package repository

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/marciomarinho/show-service/internal/domain"
)

type ShowRepository interface {
	Put(show domain.Show) error
	List() ([]domain.Show, error)
}

type DynamoInterface interface {
	GetClient() *dynamodb.Client
	GetTableName() string
}

type ShowRepo struct {
	db    *dynamodb.Client
	table string
}

func NewShowRepository(d DynamoInterface) *ShowRepo {
	return &ShowRepo{db: d.GetClient(), table: d.GetTableName()}
}

func (r *ShowRepo) Put(show domain.Show) error {
	if show.Slug == "" {
		return errors.New("slug is required")
	}

	// Derive GSI attributes for efficient querying
	var drmKey int
	if show.DRM != nil && *show.DRM {
		drmKey = 1
	} else {
		drmKey = 0
	}
	show.DRMKey = &drmKey

	// Ensure episodeCount is set for indexing
	if show.EpisodeCount == nil {
		zero := 0
		show.EpisodeCount = &zero
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
	// Use GSI to efficiently query only shows with DRM enabled and episodes > 0
	out, err := r.db.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              &r.table,
		IndexName:              aws.String("gsi_drm_episode"),
		KeyConditionExpression: aws.String("#dk = :true AND #ec > :zero"),
		ExpressionAttributeNames: map[string]string{
			"#dk":  "drmKey",
			"#ec":  "episodeCount",
			"#img": "image",
			"#si":  "showImage",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":true": &types.AttributeValueMemberN{Value: "1"},
			":zero": &types.AttributeValueMemberN{Value: "0"},
		},
		// Return only the fields needed for the response
		ProjectionExpression: aws.String("slug, title, #img.#si"),
		ScanIndexForward:     aws.Bool(true),
	})

	if err != nil {
		return nil, err
	}

	// Unmarshal into a slim struct to avoid needing all fields
	type slimShow struct {
		Slug  string        `dynamodbav:"slug"`
		Title string        `dynamodbav:"title"`
		Image *domain.Image `dynamodbav:"image"`
	}

	var rows []slimShow
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &rows); err != nil {
		return nil, err
	}

	// Map to response model with only required fields
	shows := make([]domain.Show, 0, len(rows))
	for _, row := range rows {
		shows = append(shows, domain.Show{
			Slug:  row.Slug,
			Title: row.Title,
			Image: row.Image,
		})
	}

	return shows, nil
}

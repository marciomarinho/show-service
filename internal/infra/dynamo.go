package infra

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/marciomarinho/show-service/internal/config"
	"github.com/marciomarinho/show-service/internal/repository"
)

// DynamoInterface defines the contract for DynamoDB operations
type DynamoInterface interface {
	GetClient() *dynamodb.Client
	GetTableName() string
}

// Ensure Dynamo implements repository.DynamoAPI
var _ repository.DynamoAPI = (*Dynamo)(nil)

type Dynamo struct {
	Client    *dynamodb.Client
	tableName string
}

func (d *Dynamo) GetClient() *dynamodb.Client {
	return d.Client
}

func (d *Dynamo) GetTableName() string {
	return d.tableName
}

// Implement repository.DynamoAPI methods
func (d *Dynamo) PutItem(ctx context.Context, in *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return d.Client.PutItem(ctx, in, optFns...)
}

func (d *Dynamo) Query(ctx context.Context, in *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return d.Client.Query(ctx, in, optFns...)
}

func (d *Dynamo) Scan(ctx context.Context, in *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return d.Client.Scan(ctx, in, optFns...)
}

func (d *Dynamo) TableName() string {
	return d.tableName
}

func NewDynamo(ctx context.Context, cfg *config.Config) (repository.DynamoAPI, error) {
	var ac aws.Config
	var err error

	if cfg.Env == config.EnvLocal && cfg.DynamoDB.EndpointOverride != "" {
		ac, err = awscfg.LoadDefaultConfig(ctx,
			awscfg.WithRegion(cfg.DynamoDB.Region),
			awscfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
		)
		if err != nil {
			return nil, fmt.Errorf("aws cfg local: %w", err)
		}
		ac.BaseEndpoint = aws.String(cfg.DynamoDB.EndpointOverride)
	} else {
		ac, err = awscfg.LoadDefaultConfig(ctx, awscfg.WithRegion(cfg.DynamoDB.Region))
		if err != nil {
			return nil, fmt.Errorf("aws cfg prod: %w", err)
		}
	}

	client := dynamodb.NewFromConfig(ac)

	return &Dynamo{
		Client:    client,
		tableName: cfg.DynamoDB.ShowsTable,
	}, nil
}

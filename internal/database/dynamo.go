package database

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/marciomarinho/show-service/internal/config"
)

type DynamoAPI interface {
	PutItem(ctx context.Context, in *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(ctx context.Context, in *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	Scan(ctx context.Context, in *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	TableName() string
}

type RealDynamo struct {
	Client *dynamodb.Client
	Table  string
}

func (r *RealDynamo) PutItem(ctx context.Context, in *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return r.Client.PutItem(ctx, in, optFns...)
}

func (r *RealDynamo) Query(ctx context.Context, in *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return r.Client.Query(ctx, in, optFns...)
}

func (r *RealDynamo) Scan(ctx context.Context, in *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return r.Client.Scan(ctx, in, optFns...)
}

func (r *RealDynamo) TableName() string {
	return r.Table
}

func NewDynamo(ctx context.Context, cfg *config.Config) (DynamoAPI, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

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

	return &RealDynamo{
		Client: client,
		Table:  cfg.DynamoDB.ShowsTable,
	}, nil
}

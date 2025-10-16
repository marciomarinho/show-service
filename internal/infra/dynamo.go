package infra

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/marciomarinho/show-service/internal/config"
)

// DynamoInterface defines the contract for DynamoDB operations
type DynamoInterface interface {
	GetClient() *dynamodb.Client
	GetTableName() string
}

type Dynamo struct {
	Client    *dynamodb.Client
	TableName string
}

func (d *Dynamo) GetClient() *dynamodb.Client {
	return d.Client
}

func (d *Dynamo) GetTableName() string {
	return d.TableName
}

func NewDynamo(ctx context.Context, cfg *config.Config) (DynamoInterface, error) {
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
		TableName: cfg.DynamoDB.ShowsTable,
	}, nil
}

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

type Dynamo struct {
	Client    *dynamodb.Client
	TableName string
}

func NewDynamo(ctx context.Context, cfg *config.Config) (*Dynamo, error) {
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

	// Optionally create table in local
	// if cfg.DynamoDB.CreateTableIfMissing {
	// 	if err := ensureShowsTable(ctx, client, cfg.DynamoDB.ShowsTable); err != nil {
	// 		return nil, err
	// 	}
	// }

	return &Dynamo{
		Client:    client,
		TableName: cfg.DynamoDB.ShowsTable,
	}, nil
}

// func ensureShowsTable(ctx context.Context, db *dynamodb.Client, table string) error {
// 	// Check if table exists
// 	_, err := db.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: &table})
// 	if err == nil {
// 		return nil
// 	}
// 	var rnfe *types.ResourceNotFoundException
// 	if !errors.As(err, &rnfe) {
// 		return fmt.Errorf("describe table: %w", err)
// 	}

// 	// Create minimal PK-only table; adjust as needed
// 	_, err = db.CreateTable(ctx, &dynamodb.CreateTableInput{
// 		TableName: &table,
// 		AttributeDefinitions: []types.AttributeDefinition{
// 			{AttributeName: aws.String("slug"), AttributeType: types.ScalarAttributeTypeS},
// 		},
// 		KeySchema: []types.KeySchemaElement{
// 			{AttributeName: aws.String("slug"), KeyType: types.KeyTypeHash},
// 		},
// 		BillingMode: types.BillingModePayPerRequest,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("create table: %w", err)
// 	}
// 	return nil
// }

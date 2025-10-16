package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoAPI is the interface your repo uses (easy to mock).
type DynamoAPI interface {
	PutItem(ctx context.Context, in *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(ctx context.Context, in *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	Scan(ctx context.Context, in *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	TableName() string
}

// RealDynamo adapts the AWS client to DynamoAPI.
type RealDynamo struct {
	Client   *dynamodb.Client
	Table    string
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
func (r *RealDynamo) TableName() string { return r.Table }

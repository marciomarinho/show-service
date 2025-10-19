package database

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/require"

	"github.com/marciomarinho/show-service/internal/config"
	mocks "github.com/marciomarinho/show-service/internal/database/mocks"
)

func TestRealDynamo_PutItem(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		input          *dynamodb.PutItemInput
		mockReturn     *dynamodb.PutItemOutput
		mockError      error
		expectedOutput *dynamodb.PutItemOutput
		expectedError  error
	}{
		{
			name: "successful put item",
			ctx:  context.Background(),
			input: &dynamodb.PutItemInput{
				TableName: aws.String("test-table"),
				Item: map[string]types.AttributeValue{
					"id": &types.AttributeValueMemberS{Value: "test-id"},
				},
			},
			mockReturn: &dynamodb.PutItemOutput{
				Attributes: map[string]types.AttributeValue{
					"id": &types.AttributeValueMemberS{Value: "test-id"},
				},
			},
			expectedOutput: &dynamodb.PutItemOutput{
				Attributes: map[string]types.AttributeValue{
					"id": &types.AttributeValueMemberS{Value: "test-id"},
				},
			},
		},
		{
			name: "put item with nil context",
			ctx:  nil,
			input: &dynamodb.PutItemInput{
				TableName: aws.String("test-table"),
			},
			mockReturn: &dynamodb.PutItemOutput{},
		},
		{
			name:       "put item with nil input",
			ctx:        context.Background(),
			input:      nil,
			mockReturn: &dynamodb.PutItemOutput{},
		},
		{
			name: "put item with AWS error",
			ctx:  context.Background(),
			input: &dynamodb.PutItemInput{
				TableName: aws.String("test-table"),
			},
			mockError:     errors.New("ConditionalCheckFailedException"),
			expectedError: errors.New("ConditionalCheckFailedException"),
		},
		{
			name: "put item with nil return",
			ctx:  context.Background(),
			input: &dynamodb.PutItemInput{
				TableName: aws.String("test-table"),
			},
			mockReturn: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockDynamoAPI(t)

			testAPI := &testDynamoAPI{mock: mockClient}

			if tt.mockError != nil {
				mockClient.EXPECT().PutItem(tt.ctx, tt.input).Return(nil, tt.mockError)
			} else {
				mockClient.EXPECT().PutItem(tt.ctx, tt.input).Return(tt.mockReturn, nil)
			}

			output, err := testAPI.PutItem(tt.ctx, tt.input)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			if tt.expectedOutput != nil {
				require.NotNil(t, output)
				require.Equal(t, tt.expectedOutput.Attributes, output.Attributes)
			}
		})
	}
}

func TestRealDynamo_Query(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		input          *dynamodb.QueryInput
		mockReturn     *dynamodb.QueryOutput
		mockError      error
		expectedOutput *dynamodb.QueryOutput
		expectedError  error
	}{
		{
			name: "successful query",
			ctx:  context.Background(),
			input: &dynamodb.QueryInput{
				TableName: aws.String("test-table"),
			},
			mockReturn: &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{
					{
						"id": &types.AttributeValueMemberS{Value: "test-id"},
					},
				},
				Count: 1,
			},
			expectedOutput: &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{
					{
						"id": &types.AttributeValueMemberS{Value: "test-id"},
					},
				},
				Count: 1,
			},
		},
		{
			name: "query with nil context",
			ctx:  nil,
			input: &dynamodb.QueryInput{
				TableName: aws.String("test-table"),
			},
			mockReturn: &dynamodb.QueryOutput{},
		},
		{
			name:       "query with nil input",
			ctx:        context.Background(),
			input:      nil,
			mockReturn: &dynamodb.QueryOutput{},
		},
		{
			name: "query with AWS error",
			ctx:  context.Background(),
			input: &dynamodb.QueryInput{
				TableName: aws.String("test-table"),
			},
			mockError:     errors.New("ValidationException"),
			expectedError: errors.New("ValidationException"),
		},
		{
			name: "query with empty results",
			ctx:  context.Background(),
			input: &dynamodb.QueryInput{
				TableName: aws.String("test-table"),
			},
			mockReturn: &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{},
				Count: 0,
			},
			expectedOutput: &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{},
				Count: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockDynamoAPI(t)

			testAPI := &testDynamoAPI{mock: mockClient}

			if tt.mockError != nil {
				mockClient.EXPECT().Query(tt.ctx, tt.input).Return(nil, tt.mockError)
			} else {
				mockClient.EXPECT().Query(tt.ctx, tt.input).Return(tt.mockReturn, nil)
			}

			output, err := testAPI.Query(tt.ctx, tt.input)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			if tt.expectedOutput != nil {
				require.NotNil(t, output)
				require.Equal(t, tt.expectedOutput.Count, output.Count)
				require.Equal(t, tt.expectedOutput.Items, output.Items)
			}
		})
	}
}

func TestRealDynamo_Scan(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		input          *dynamodb.ScanInput
		mockReturn     *dynamodb.ScanOutput
		mockError      error
		expectedOutput *dynamodb.ScanOutput
		expectedError  error
	}{
		{
			name: "successful scan",
			ctx:  context.Background(),
			input: &dynamodb.ScanInput{
				TableName: aws.String("test-table"),
			},
			mockReturn: &dynamodb.ScanOutput{
				Items: []map[string]types.AttributeValue{
					{
						"id": &types.AttributeValueMemberS{Value: "test-id"},
					},
				},
				Count: 1,
			},
			expectedOutput: &dynamodb.ScanOutput{
				Items: []map[string]types.AttributeValue{
					{
						"id": &types.AttributeValueMemberS{Value: "test-id"},
					},
				},
				Count: 1,
			},
		},
		{
			name: "scan with nil context",
			ctx:  nil,
			input: &dynamodb.ScanInput{
				TableName: aws.String("test-table"),
			},
			mockReturn: &dynamodb.ScanOutput{},
		},
		{
			name:       "scan with nil input",
			ctx:        context.Background(),
			input:      nil,
			mockReturn: &dynamodb.ScanOutput{},
		},
		{
			name: "scan with AWS error",
			ctx:  context.Background(),
			input: &dynamodb.ScanInput{
				TableName: aws.String("test-table"),
			},
			mockError:     errors.New("ResourceNotFoundException"),
			expectedError: errors.New("ResourceNotFoundException"),
		},
		{
			name: "scan with empty results",
			ctx:  context.Background(),
			input: &dynamodb.ScanInput{
				TableName: aws.String("test-table"),
			},
			mockReturn: &dynamodb.ScanOutput{
				Items: []map[string]types.AttributeValue{},
				Count: 0,
			},
			expectedOutput: &dynamodb.ScanOutput{
				Items: []map[string]types.AttributeValue{},
				Count: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockDynamoAPI(t)

			testAPI := &testDynamoAPI{mock: mockClient}

			if tt.mockError != nil {
				mockClient.EXPECT().Scan(tt.ctx, tt.input).Return(nil, tt.mockError)
			} else {
				mockClient.EXPECT().Scan(tt.ctx, tt.input).Return(tt.mockReturn, nil)
			}

			output, err := testAPI.Scan(tt.ctx, tt.input)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			if tt.expectedOutput != nil {
				require.NotNil(t, output)
				require.Equal(t, tt.expectedOutput.Count, output.Count)
				require.Equal(t, tt.expectedOutput.Items, output.Items)
			}
		})
	}
}

func TestRealDynamo_TableName(t *testing.T) {
	tests := []struct {
		name           string
		table          string
		expectedOutput string
	}{
		{
			name:           "returns configured table name",
			table:          "shows-local",
			expectedOutput: "shows-local",
		},
		{
			name:           "returns empty table name",
			table:          "",
			expectedOutput: "",
		},
		{
			name:           "returns custom table name",
			table:          "custom-table",
			expectedOutput: "custom-table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAPI := &testDynamoAPI{tableName: tt.table}

			output := testAPI.TableName()

			require.Equal(t, tt.expectedOutput, output)
		})
	}
}

type testDynamoAPI struct {
	mock      *mocks.MockDynamoAPI
	tableName string
}

func (t *testDynamoAPI) PutItem(ctx context.Context, in *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if t.mock != nil {
		return t.mock.PutItem(ctx, in, optFns...)
	}
	return &dynamodb.PutItemOutput{}, nil
}

func (t *testDynamoAPI) Query(ctx context.Context, in *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	if t.mock != nil {
		return t.mock.Query(ctx, in, optFns...)
	}
	return &dynamodb.QueryOutput{}, nil
}

func (t *testDynamoAPI) Scan(ctx context.Context, in *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	if t.mock != nil {
		return t.mock.Scan(ctx, in, optFns...)
	}
	return &dynamodb.ScanOutput{}, nil
}

func (t *testDynamoAPI) TableName() string {
	return t.tableName
}

func TestNewDynamo(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		config      *config.Config
		expectError bool
		errorMsg    string
		validateFn  func(t *testing.T, result DynamoAPI)
	}{
		{
			name: "successful local configuration with endpoint",
			ctx:  context.Background(),
			config: &config.Config{
				Env: config.EnvLocal,
				DynamoDB: config.DynamoDB{
					Region:           "us-east-1",
					EndpointOverride: "http://localhost:8000",
					ShowsTable:       "shows-local",
				},
			},
			expectError: false,
			validateFn: func(t *testing.T, result DynamoAPI) {
				require.NotNil(t, result)
				rd, ok := result.(*RealDynamo)
				require.True(t, ok, "Expected *RealDynamo")
				require.Equal(t, "shows-local", rd.TableName())
			},
		},
		{
			name: "successful production configuration",
			ctx:  context.Background(),
			config: &config.Config{
				Env: config.EnvDev,
				DynamoDB: config.DynamoDB{
					Region:     "ap-southeast-2",
					ShowsTable: "shows-dev",
				},
			},
			expectError: false,
			validateFn: func(t *testing.T, result DynamoAPI) {
				require.NotNil(t, result)
				rd, ok := result.(*RealDynamo)
				require.True(t, ok, "Expected *RealDynamo")
				require.Equal(t, "shows-dev", rd.TableName())
			},
		},
		{
			name: "local configuration without endpoint",
			ctx:  context.Background(),
			config: &config.Config{
				Env: config.EnvLocal,
				DynamoDB: config.DynamoDB{
					Region:     "us-east-1",
					ShowsTable: "shows-local",
				},
			},
			expectError: false,
			validateFn: func(t *testing.T, result DynamoAPI) {
				require.NotNil(t, result)
				rd, ok := result.(*RealDynamo)
				require.True(t, ok, "Expected *RealDynamo")
				require.Equal(t, "shows-local", rd.TableName())
			},
		},
		{
			name: "nil context",
			ctx:  nil,
			config: &config.Config{
				Env: config.EnvLocal,
				DynamoDB: config.DynamoDB{
					Region:     "us-east-1",
					ShowsTable: "shows-local",
				},
			},
			expectError: true,
			errorMsg:    "context is required",
		},
		{
			name:        "nil config",
			ctx:         context.Background(),
			config:      nil,
			expectError: true,
			errorMsg:    "config is required",
		},
		{
			name: "config with empty region",
			ctx:  context.Background(),
			config: &config.Config{
				Env: config.EnvLocal,
				DynamoDB: config.DynamoDB{
					Region:     "",
					ShowsTable: "shows-local",
				},
			},
			expectError: false,
			validateFn: func(t *testing.T, result DynamoAPI) {
				require.NotNil(t, result)
			},
		},
		{
			name: "config with empty table name",
			ctx:  context.Background(),
			config: &config.Config{
				Env: config.EnvLocal,
				DynamoDB: config.DynamoDB{
					Region:     "us-east-1",
					ShowsTable: "",
				},
			},
			expectError: false,
			validateFn: func(t *testing.T, result DynamoAPI) {
				require.NotNil(t, result)
				rd, ok := result.(*RealDynamo)
				require.True(t, ok, "Expected *RealDynamo")
				require.Equal(t, "", rd.TableName())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewDynamo(tt.ctx, tt.config)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				if tt.validateFn != nil {
					tt.validateFn(t, result)
				}
			}
		})
	}
}

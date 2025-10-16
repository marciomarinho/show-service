#!/bin/bash
set -e

export AWS_ACCESS_KEY_ID=dummy
export AWS_SECRET_ACCESS_KEY=dummy
export AWS_DEFAULT_REGION=ap-southeast-2

# Start DynamoDB Local in the background
java -jar DynamoDBLocal.jar -inMemory -port 8000 &

# Wait for DynamoDB to be ready
echo "Waiting for DynamoDB Local to start..."
sleep 5

# Create the shows table with the correct structure for the application
echo "Creating DynamoDB table with GSI..."
aws dynamodb create-table \
  --table-name shows-local \
  --attribute-definitions \
    AttributeName=slug,AttributeType=S \
    AttributeName=drmKey,AttributeType=N \
    AttributeName=episodeCount,AttributeType=N \
  --key-schema AttributeName=slug,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --global-secondary-indexes \
    "[{\"IndexName\": \"gsi_drm_episode\",\"KeySchema\": [{\"AttributeName\": \"drmKey\",\"KeyType\": \"HASH\"},{\"AttributeName\": \"episodeCount\",\"KeyType\": \"RANGE\"}],\"Projection\": {\"ProjectionType\": \"ALL\"}}]" \
  --endpoint-url http://localhost:8000 \
  --region ap-southeast-2

if [ $? -eq 0 ]; then
    echo "Table 'shows-local' created successfully with GSI!"
else
    echo "Failed to create table"
    exit 1
fi

# Keep the container running
wait
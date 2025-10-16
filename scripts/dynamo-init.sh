#!/usr/bin/env sh
set -eu

ENDPOINT_URL="${ENDPOINT_URL:-http://localhost:8000}"
TABLE_NAME="${TABLE_NAME:-shows_local}"
AWS_REGION="${AWS_REGION:-local}"

echo "Waiting for DynamoDB on ${ENDPOINT_URL}..."
for i in $(seq 1 60); do
  if aws dynamodb list-tables --endpoint-url "${ENDPOINT_URL}" --region "${AWS_REGION}" >/dev/null 2>&1; then
    echo "DynamoDB is up."
    break
  fi
  sleep 1
done

# Create table if missing
if aws dynamodb describe-table \
    --table-name "${TABLE_NAME}" \
    --endpoint-url "${ENDPOINT_URL}" \
    --region "${AWS_REGION}" >/dev/null 2>&1; then
  echo "Table '${TABLE_NAME}' already exists."
else
  echo "Creating table '${TABLE_NAME}' ..."
  aws dynamodb create-table \
    --table-name "${TABLE_NAME}" \
    --attribute-definitions \
      AttributeName=slug,AttributeType=S \
      AttributeName=drmKey,AttributeType=N \
      AttributeName=episodeCount,AttributeType=N \
    --key-schema AttributeName=slug,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --table-class STANDARD \
    --global-secondary-indexes \
      "[{\"IndexName\": \"gsi_drm_episode\",\"KeySchema\": [{\"AttributeName\": \"drmKey\",\"KeyType\": \"HASH\"},{\"AttributeName\": \"episodeCount\",\"KeyType\": \"RANGE\"}],\"Projection\": {\"ProjectionType\": \"ALL\"}}]" \
    --endpoint-url "${ENDPOINT_URL}" \
    --region "${AWS_REGION}"

  echo "Waiting for table to exist ..."
  aws dynamodb wait table-exists \
    --table-name "${TABLE_NAME}" \
    --endpoint-url "${ENDPOINT_URL}" \
    --region "${AWS_REGION}"
fi

# Optional: seed one item if present
if [ -f "/seed/seed-one-item.json" ]; then
  echo "Seeding sample item..."
  aws dynamodb put-item \
    --table-name "${TABLE_NAME}" \
    --item file:///seed/seed-one-item.json \
    --endpoint-url "${ENDPOINT_URL}" \
    --region "${AWS_REGION}"
fi

echo "Init complete."

#!/bin/bash

# Simple curl command to test GitHub Meta API
# Usage: ./curl-github-meta.sh

API_URL="http://localhost:8080/api/v1/fetch-json"
GITHUB_META_URL="https://api.github.com/meta"

echo "Calling API with GitHub Meta URL..."
echo "API Endpoint: $API_URL"
echo "Parameter: $GITHUB_META_URL"
echo ""

# Make the curl request
curl -X GET "${API_URL}?link=${GITHUB_META_URL}" \
  -H "Accept: application/json" \
  -H "Content-Type: application/json" \
  --verbose \
  | jq .

echo ""
echo "Done!"
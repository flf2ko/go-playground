#!/bin/bash

# Test script for Go API Sample with GitHub Meta API
# Usage: ./test-github-api.sh

set -e

API_BASE_URL="http://localhost:8080"
GITHUB_META_URL="https://api.github.com/meta"

echo "🚀 Testing Go API Sample with GitHub Meta API"
echo "================================================"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to check if API is running
check_api_health() {
    echo -e "${YELLOW}1. Checking API health...${NC}"
    
    if curl -s "${API_BASE_URL}/health" > /dev/null; then
        echo -e "${GREEN}✅ API is running${NC}"
        curl -s "${API_BASE_URL}/health" | jq .
    else
        echo -e "${RED}❌ API is not running. Please start it first with: make dev-up${NC}"
        exit 1
    fi
    echo ""
}

# Function to fetch GitHub meta data
fetch_github_meta() {
    echo -e "${YELLOW}2. Fetching GitHub Meta API data...${NC}"
    echo "URL: ${GITHUB_META_URL}"
    
    # URL encode the GitHub meta URL
    ENCODED_URL=$(python3 -c "import urllib.parse; print(urllib.parse.quote('${GITHUB_META_URL}', safe=''))")
    
    echo "Making request to: ${API_BASE_URL}/api/v1/fetch-json?link=${GITHUB_META_URL}"
    
    RESPONSE=$(curl -s "${API_BASE_URL}/api/v1/fetch-json?link=${GITHUB_META_URL}")
    
    if echo "$RESPONSE" | jq -e '.success == true' > /dev/null; then
        echo -e "${GREEN}✅ Successfully fetched and stored GitHub meta data${NC}"
        echo "$RESPONSE" | jq .
    else
        echo -e "${RED}❌ Failed to fetch GitHub meta data${NC}"
        echo "$RESPONSE" | jq .
        exit 1
    fi
    echo ""
}

# Function to get stored records
get_stored_records() {
    echo -e "${YELLOW}3. Retrieving stored records...${NC}"
    
    RECORDS=$(curl -s "${API_BASE_URL}/api/v1/records")
    
    if echo "$RECORDS" | jq -e '.success == true' > /dev/null; then
        echo -e "${GREEN}✅ Successfully retrieved stored records${NC}"
        echo "$RECORDS" | jq .
        
        # Show count
        RECORD_COUNT=$(echo "$RECORDS" | jq '.count')
        echo -e "${GREEN}📊 Total records: ${RECORD_COUNT}${NC}"
    else
        echo -e "${RED}❌ Failed to retrieve records${NC}"
        echo "$RECORDS" | jq .
    fi
    echo ""
}

# Function to show the stored GitHub meta content
show_github_meta_content() {
    echo -e "${YELLOW}4. Showing GitHub Meta content preview...${NC}"
    
    RECORDS=$(curl -s "${API_BASE_URL}/api/v1/records")
    
    # Find the GitHub meta record and show its content
    GITHUB_RECORD=$(echo "$RECORDS" | jq -r '.data[] | select(.url == "https://api.github.com/meta") | .content' | head -1)
    
    if [ "$GITHUB_RECORD" != "null" ] && [ "$GITHUB_RECORD" != "" ]; then
        echo -e "${GREEN}✅ GitHub Meta API content:${NC}"
        echo "$GITHUB_RECORD" | jq .
    else
        echo -e "${YELLOW}⚠️  No GitHub Meta record found in stored data${NC}"
    fi
    echo ""
}

# Main execution
main() {
    echo "Starting API test at $(date)"
    echo ""
    
    # Check if jq is installed
    if ! command -v jq &> /dev/null; then
        echo -e "${RED}❌ jq is required but not installed. Please install it first:${NC}"
        echo "  macOS: brew install jq"
        echo "  Ubuntu: sudo apt-get install jq"
        exit 1
    fi
    
    # Run tests
    check_api_health
    fetch_github_meta
    get_stored_records
    show_github_meta_content
    
    echo -e "${GREEN}🎉 Test completed successfully!${NC}"
    echo ""
    echo "💡 You can also test manually with:"
    echo "   curl \"${API_BASE_URL}/api/v1/fetch-json?link=${GITHUB_META_URL}\""
}

# Run main function
main "$@"
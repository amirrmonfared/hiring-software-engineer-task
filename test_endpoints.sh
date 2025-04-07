#!/usr/bin/env bash

set -e

BASE_URL="http://localhost:8080"

echo "----- [1] Health Check -----"
curl -sS -X GET "$BASE_URL/health" | jq .
echo

echo "----- [2] Create Line Item -----"
CREATE_RESPONSE=$(curl -sS -X POST "$BASE_URL/api/v1/lineitems" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Banner",
    "advertiser_id": "adv123",
    "bid": 2.5,
    "budget": 1000.0,
    "placement": "homepage_top",
    "categories": ["electronics", "sale"],
    "keywords": ["summer", "discount"]
  }')
echo "$CREATE_RESPONSE" | jq .
echo

LINE_ITEM_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id')

echo "Created line item with ID: $LINE_ITEM_ID"
echo

echo "----- [3] Get All Line Items -----"
curl -sS -X GET "$BASE_URL/api/v1/lineitems" | jq .
echo

echo "----- [4] Get Line Item by ID -----"
curl -sS -X GET "$BASE_URL/api/v1/lineitems/$LINE_ITEM_ID" | jq .
echo

echo "----- [5] Get Winning Ads -----"
curl -sS -X GET "$BASE_URL/api/v1/ads?placement=homepage_top&category=electronics&keyword=discount&limit=1" | jq .
echo

echo "----- [6] Track an Impression -----"
curl -sS -X POST "$BASE_URL/api/v1/tracking" \
  -H "Content-Type: application/json" \
  -d "{
    \"event_type\": \"impression\",
    \"line_item_id\": \"$LINE_ITEM_ID\",
    \"user_id\": \"test_user\",
    \"metadata\": { \"device\": \"mobile\" }
  }" | jq .
echo

echo "All tests completed successfully."

#!/bin/bash

BASE_URL="http://localhost:8080"

# Function to print response and check for errors
check_response() {
    if [ $? -ne 0 ]; then
        echo "Error: Failed to connect to the server"
        exit 1
    fi
    if [[ $1 == *"error"* ]]; then
        echo "Error in response: $1"
    else
        echo "Success: $1"
    fi
}

echo "Testing API endpoints..."

# List Users
echo "Listing Users:"
response=$(curl -s $BASE_URL/users)
check_response "$response"

# Create a User
echo "Creating a User:"
response=$(curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123", "provider": "local"}')
check_response "$response"

# Extract user ID from response (assuming it's returned)
user_id=$(echo $response | grep -o '"id":[0-9]*' | grep -o '[0-9]*')

if [ -z "$user_id" ]; then
    echo "Failed to get user ID, using 1 for subsequent tests"
    user_id=1
fi

# Get a Specific User
echo "Getting User $user_id:"
response=$(curl -s $BASE_URL/users/$user_id)
check_response "$response"

# Create an Account
echo "Creating an Account:"
response=$(curl -s -X POST $BASE_URL/accounts \
  -H "Content-Type: application/json" \
  -d "{\"userId\": $user_id, \"accountType\": \"savings\"}")
check_response "$response"

# List Accounts
echo "Listing Accounts:"
response=$(curl -s $BASE_URL/accounts)
check_response "$response"

# Create an Order
echo "Creating an Order:"
response=$(curl -s -X POST $BASE_URL/orders \
  -H "Content-Type: application/json" \
  -d "{
    \"accountId\": $user_id,
    \"assetId\": 1,
    \"orderType\": \"market\",
    \"side\": \"buy\",
    \"quantity\": 10,
    \"price\": 100.50
  }")
check_response "$response"

# List Orders
echo "Listing Orders:"
response=$(curl -s $BASE_URL/orders)
check_response "$response"

# Extract order ID from response (assuming it's returned)
order_id=$(echo $response | grep -o '"id":[0-9]*' | grep -o '[0-9]*' | head -n1)

if [ -z "$order_id" ]; then
    echo "Failed to get order ID, using 1 for subsequent tests"
    order_id=1
fi

# Cancel an Order
echo "Cancelling Order $order_id:"
response=$(curl -s -X DELETE $BASE_URL/orders/$order_id)
check_response "$response"

# Create a Portfolio
echo "Creating a Portfolio:"
response=$(curl -s -X POST $BASE_URL/portfolios \
  -H "Content-Type: application/json" \
  -d "{\"accountId\": $user_id}")
check_response "$response"

echo "API testing completed."
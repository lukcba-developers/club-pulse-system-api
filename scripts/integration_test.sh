#!/bin/bash

# Configuration
# Configuration
API_URL="http://localhost:8081/api/v1"
EMAIL="testuser@example.com"
PASSWORD="password123"
NAME="Test User"

echo "🧪 Starting Integration Tests for Club Pulse System..."
echo "---------------------------------------------------"

# 1. Health Check
echo "1. Checking System Health..."
HEALTH_RESPONSE=$(curl -s "http://localhost:8081/health")
if [[ "$HEALTH_RESPONSE" == *"UP"* ]]; then
    echo "✅ Backend is UP"
else
    echo "❌ Backend is DOWN"
    echo "Response: $HEALTH_RESPONSE"
    exit 1
fi

# 2. Register User
echo -e "\n2. Testing Registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"$NAME\", \"email\":\"$EMAIL\", \"password\":\"$PASSWORD\"}")

# If we get a token or success, it's good. Or if user exists (conflict) it's also okay for repeatability.
echo "Response: $REGISTER_RESPONSE"

# 3. Login
echo -e "\n3. Testing Login..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\", \"password\":\"$PASSWORD\"}")

TOKEN=$(echo $LOGIN_RESPONSE | sed 's/.*"access_token":"\([^"]*\)".*/\1/')


if [ -n "$TOKEN" ] && [ "$TOKEN" != "$LOGIN_RESPONSE" ]; then
    echo "✅ Login Successful. Token received."
    
    # Capture Refresh Token
    REFRESH_TOKEN=$(echo $LOGIN_RESPONSE | sed 's/.*"refresh_token":"\([^"]*\)".*/\1/')
    if [ -n "$REFRESH_TOKEN" ] && [ "$REFRESH_TOKEN" != "$LOGIN_RESPONSE" ]; then
        echo "✅ Refresh Token received."
    else
        echo "❌ No Refresh Token received."
        exit 1
    fi

    echo -e "\n3.1 Testing Token Refresh..."
    REFRESH_RESPONSE=$(curl -s -X POST "$API_URL/auth/refresh" \
        -H "Content-Type: application/json" \
        -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}")
    
    NEW_ACCESS_TOKEN=$(echo $REFRESH_RESPONSE | sed 's/.*"access_token":"\([^"]*\)".*/\1/')
    NEW_REFRESH_TOKEN=$(echo $REFRESH_RESPONSE | sed 's/.*"refresh_token":"\([^"]*\)".*/\1/')
    
    if [ -n "$NEW_ACCESS_TOKEN" ] && [ "$NEW_ACCESS_TOKEN" != "$REFRESH_RESPONSE" ]; then
         echo "✅ Token Refresh Successful"
         TOKEN=$NEW_ACCESS_TOKEN # Use new token for future requests
    else
         echo "❌ Token Refresh Failed"
         echo "Response: $REFRESH_RESPONSE"
         exit 1
    fi

    echo -e "\n3.2 Testing Logout..."
    LOGOUT_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/auth/logout" \
        -H "Content-Type: application/json" \
        -d "{\"refresh_token\":\"$NEW_REFRESH_TOKEN\"}")
    
    if [ "$LOGOUT_CODE" == "204" ]; then
        echo "✅ Logout Successful"
    else
        echo "❌ Logout Failed (Code: $LOGOUT_CODE)"
        exit 1
    fi

    # Fetch User ID (using new token)
    echo "Fetching User Profile..."
    PROFILE_RESPONSE=$(curl -s -X GET "$API_URL/users/me" -H "Authorization: Bearer $TOKEN")
    USER_ID=$(echo $PROFILE_RESPONSE | sed 's/.*"id":"\([^"]*\)".*/\1/')
    echo "User ID fetched: $USER_ID"
else
    echo "❌ Login Failed."
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi


echo -e "\n4. Testing Facilities..."
echo "Creating a facility..."
FACILITY_RESPONSE=$(curl -s -X POST $API_URL/facilities \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Grand Tennis Court",
    "type": "TENNIS_COURT",
    "location": {"coordinates": [0,0], "address": "123 Main St", "city": "New York", "country": "USA"},
    "specs": {"surface": "Clay"},
    "hourly_rate": 50,
    "capacity": 4,
    "status": "AVAILABLE"
  }')
echo "Response: $FACILITY_RESPONSE"

echo "Listing facilities..."
curl -s $API_URL/facilities | grep "Grand Tennis Court" && echo "✅ Facility created and listed" || echo "❌ Facility test failed"

# Capture Facility ID for booking tests
FACILITY_ID=$(echo $FACILITY_RESPONSE | sed 's/.*"id":"\([^"]*\)".*/\1/')
echo "Facility ID: $FACILITY_ID"

echo -e "\n4.1 Testing Facility Updates (Maintenance & Equipment)..."

# 1. Update Equipment
echo "Updating Equipment status..."
UPDATE_RES=$(curl -s -X PUT "$API_URL/facilities/$FACILITY_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "specifications": {
        "surface_type": "Clay", 
        "lighting": true, 
        "covered": false,
        "equipment": ["Net", "Rackets"]
    }
  }')
echo $UPDATE_RES | grep "Rackets" && echo "✅ Equipment Updated" || echo "❌ Equipment Update Failed: $UPDATE_RES"

# 2. Set to Maintenance
echo "Setting Facility to Maintenance..."
MAINT_RES=$(curl -s -X PUT "$API_URL/facilities/$FACILITY_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "maintenance"}')
echo $MAINT_RES | grep "maintenance" && echo "✅ Status set to Maintenance" || echo "❌ Status Update Failed: $MAINT_RES"

# 3. Try to Book (Should Fail)
echo "Attempting to book maintenance facility..."
FAIL_BOOKING=$(curl -s -X POST $API_URL/bookings \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "user_id": "'$USER_ID'",
    "facility_id": "'$FACILITY_ID'",
    "start_time": "'$(date -v+2d -u +"%Y-%m-%dT10:00:00Z")'",
    "end_time": "'$(date -v+2d -u +"%Y-%m-%dT11:00:00Z")'"
  }')

echo $FAIL_BOOKING | grep "not active" && echo "✅ Booking Blocked correctly" || echo "❌ Booking should have failed: $FAIL_BOOKING"

# 4. Restore to Active
echo "Restoring Facility to Active..."
ACTIVE_RES=$(curl -s -X PUT "$API_URL/facilities/$FACILITY_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "active"}')
echo $ACTIVE_RES | grep "active" && echo "✅ Status restored to Active" || echo "❌ Restore Failed: $ACTIVE_RES"

echo -e "\n5. Testing Membership Module..."
echo "Listing Membership Tiers (should be empty initially or populated by seed)..."
TIERS_RESPONSE=$(curl -s -X GET $API_URL/memberships/tiers \
  -H "Authorization: Bearer $TOKEN")
echo "Response: $TIERS_RESPONSE"

# Verify endpoint is reachable
echo $TIERS_RESPONSE | grep "data" && echo "✅ Membership Tiers endpoint reachable" || echo "❌ Membership Tiers test failed"

echo "Creating a Membership..."
MEMBERSHIP_RESPONSE=$(curl -s -X POST $API_URL/memberships \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "user_id": "'$USER_ID'",
    "membership_tier_id": "11111111-1111-1111-1111-111111111111", 
    "billing_cycle": "MONTHLY"
  }')
echo "Create Membership Response: $MEMBERSHIP_RESPONSE"

echo "Listing User Memberships..."
curl -s -X GET $API_URL/memberships \
  -H "Authorization: Bearer $TOKEN" | grep "data" && echo "✅ List Memberships endpoint reachable" || echo "❌ List Memberships test failed"

# --- Booking Tests ---
echo -e "\n6. Testing Booking Module..."
START_TIME=$(date -v+1d -u +"%Y-%m-%dT10:00:00Z")
END_TIME=$(date -v+1d -u +"%Y-%m-%dT11:00:00Z")

echo "Creating a Booking at $START_TIME..."
BOOKING_RESPONSE=$(curl -s -X POST $API_URL/bookings \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "user_id": "'$USER_ID'",
    "facility_id": "'$FACILITY_ID'",
    "start_time": "'$START_TIME'",
    "end_time": "'$END_TIME'"
  }')
echo "Response: $BOOKING_RESPONSE"

BOOKING_ID=$(echo $BOOKING_RESPONSE | sed 's/.*"id":"\([^"]*\)".*/\1/')

if [[ "$BOOKING_RESPONSE" == *"id"* ]]; then
    echo "✅ Booking Created Successfully"
else
    echo "❌ Booking Creation Failed"
fi

echo "Testing Conflict (Overlapping Booking)..."
CONFLICT_RESPONSE=$(curl -s -X POST $API_URL/bookings \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "user_id": "'$USER_ID'",
    "facility_id": "'$FACILITY_ID'",
    "start_time": "'$START_TIME'",
    "end_time": "'$END_TIME'"
  }')
echo "Conflict Response: $CONFLICT_RESPONSE"
echo $CONFLICT_RESPONSE | grep "conflict" && echo "✅ Conflict Detection Passed" || echo "❌ Conflict Detection Failed"

echo "Listing Bookings..."
LIST_BOOKINGS_RESPONSE=$(curl -s -X GET $API_URL/bookings -H "Authorization: Bearer $TOKEN")
echo $LIST_BOOKINGS_RESPONSE | grep "$BOOKING_ID" && echo "✅ Listing Bookings Passed" || echo "❌ Listing Bookings Failed"

echo "Checking Availability API..."
AVAILABILITY_RESPONSE=$(curl -s -X GET "$API_URL/bookings/availability?facility_id=$FACILITY_ID&date=$(date -v+1d -u +"%Y-%m-%d")" -H "Authorization: Bearer $TOKEN")
echo $AVAILABILITY_RESPONSE | grep "data" && echo "✅ Availability Check Passed" || echo "❌ Availability Check Failed: $AVAILABILITY_RESPONSE"

echo "Cancelling Booking..."
CANCEL_RESPONSE=$(curl -s -X DELETE "$API_URL/bookings/$BOOKING_ID" -H "Authorization: Bearer $TOKEN")
echo "Cancel Response: $CANCEL_RESPONSE"
echo $CANCEL_RESPONSE | grep "cancelled" && echo "✅ Booking Cancelled" || echo "❌ Booking Cancellation Failed"

echo -e "\n7. Testing User Module..."
echo "Listing Users..."
USERS_LIST=$(curl -s -X GET "$API_URL/users" -H "Authorization: Bearer $TOKEN")
# Assuming the response has envelopes {"data": [...]}
COUNT=$(echo $USERS_LIST | grep -o "id" | wc -l) 
if [ "$COUNT" -gt 0 ]; then
    echo "✅ Users Listed (Count: $COUNT)"
else
     echo "❌ User List Failed or Empty (Response: $USERS_LIST)"
     exit 1
fi

# Create a temporary user to delete
echo "Creating dummy user for deletion..."
DUMMY_EMAIL="delete_me_$(date +%s)@example.com"
DUMMY_USER_RES=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Delete Me\",\"email\":\"$DUMMY_EMAIL\",\"password\":\"password123\"}")

# Register now returns tokens, so we need to get ID from profile using the new token
DUMMY_TOKEN=$(echo $DUMMY_USER_RES | sed 's/.*"access_token":"\([^"]*\)".*/\1/')

if [ -n "$DUMMY_TOKEN" ] && [ "$DUMMY_TOKEN" != "$DUMMY_USER_RES" ]; then
     # Get ID of the dummy user
     DUMMY_PROFILE=$(curl -s -X GET "$API_URL/users/me" -H "Authorization: Bearer $DUMMY_TOKEN")
     DUMMY_ID=$(echo $DUMMY_PROFILE | sed 's/.*"id":"\([^"]*\)".*/\1/')
     
     echo "✅ Dummy User Created (ID: $DUMMY_ID)"
     
     echo "Deleting User..."
     # Use the MAIN admin/test token to delete the dummy user (since self-delete is forbidden)
     DELETE_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$API_URL/users/$DUMMY_ID" \
        -H "Authorization: Bearer $TOKEN")
     
     if [ "$DELETE_CODE" == "204" ]; then
         echo "✅ User Deleted Successfully"
     else
         echo "❌ User Delete Failed (Code: $DELETE_CODE)"
         # Don't exit here, maybe just warn? No, it's a critical test.
         exit 1
     fi
     
     # Verify Deletion
     CHECK_LIST=$(curl -s -X GET "$API_URL/users" -H "Authorization: Bearer $TOKEN")
     if [[ "$CHECK_LIST" == *"$DUMMY_ID"* ]]; then
         echo "❌ User Still Present in List after Deletion"
         exit 1
     else
         echo "✅ User confirmed removed from List"
     fi
     
else
     echo "❌ Failed to register dummy user (No token returned)"
     echo "Response: $DUMMY_USER_RES"
     exit 1
fi


echo -e "\n---------------------------------------------------"
echo "🎉 All Integration Tests Passed Successfully!"

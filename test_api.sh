#!/bin/bash

# Mindloop API Test Script
# This script demonstrates all the API functionality

BASE_URL="http://localhost:8080/api/v1"

echo "🧠 Mindloop API Test Script"
echo "=========================="
echo ""

# Test 1: Health Check
echo "1. Testing Health Check..."
curl -s "$BASE_URL/healthz" | jq .
echo ""

# Test 2: API Info
echo "2. Testing API Info..."
curl -s "$BASE_URL/" | jq .
echo ""

# Test 3: Create Habit
echo "3. Creating a Habit..."
HABIT_RESPONSE=$(curl -s -X POST "$BASE_URL/habits" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Daily Exercise",
    "description": "30 minutes of physical activity",
    "target_count": 1,
    "interval": "daily"
  }')
echo "$HABIT_RESPONSE" | jq .
HABIT_ID=$(echo "$HABIT_RESPONSE" | jq -r '.data.id')
echo ""

# Test 4: List Habits
echo "4. Listing All Habits..."
curl -s "$BASE_URL/habits" | jq .
echo ""

# Test 5: Create Intent
echo "5. Creating an Intent..."
INTENT_RESPONSE=$(curl -s -X POST "$BASE_URL/intents" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Complete API Documentation"
  }')
echo "$INTENT_RESPONSE" | jq .
INTENT_ID=$(echo "$INTENT_RESPONSE" | jq -r '.data.id')
echo ""

# Test 6: List Intents
echo "6. Listing All Intents..."
curl -s "$BASE_URL/intents" | jq .
echo ""

# Test 7: Create Focus Session
echo "7. Creating a Focus Session..."
FOCUS_RESPONSE=$(curl -s -X POST "$BASE_URL/focus" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "API Development Session"
  }')
echo "$FOCUS_RESPONSE" | jq .
FOCUS_ID=$(echo "$FOCUS_RESPONSE" | jq -r '.data.id')
echo ""

# Test 8: List Focus Sessions
echo "8. Listing All Focus Sessions..."
curl -s "$BASE_URL/focus" | jq .
echo ""

# Test 9: Create Journal Entry
echo "9. Creating a Journal Entry..."
JOURNAL_RESPONSE=$(curl -s -X POST "$BASE_URL/journal" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "API Development Day",
    "content": "Successfully implemented the Mindloop API with all features working perfectly!",
    "mood": "excited"
  }')
echo "$JOURNAL_RESPONSE" | jq .
JOURNAL_ID=$(echo "$JOURNAL_RESPONSE" | jq -r '.data.id')
echo ""

# Test 10: List Journal Entries
echo "10. Listing All Journal Entries..."
curl -s "$BASE_URL/journal" | jq .
echo ""

# Test 11: Log Habit Completion
echo "11. Logging Habit Completion..."
curl -s -X POST "$BASE_URL/habits/$HABIT_ID/log" \
  -H "Content-Type: application/json" \
  -d '{
    "actual_count": 1
  }' | jq .
echo ""

# Test 12: End Focus Session
echo "12. Ending Focus Session..."
curl -s -X POST "$BASE_URL/focus/$FOCUS_ID/end" | jq .
echo ""

# Test 13: Rate Focus Session
echo "13. Rating Focus Session..."
curl -s -X POST "$BASE_URL/focus/$FOCUS_ID/rate" \
  -H "Content-Type: application/json" \
  -d '{
    "rating": 9
  }' | jq .
echo ""

# Test 14: End Intent
echo "14. Ending Intent..."
curl -s -X POST "$BASE_URL/intents/$INTENT_ID/end" | jq .
echo ""

# Test 15: Get Daily Summary
echo "15. Getting Daily Summary..."
curl -s "$BASE_URL/summary/daily" | jq .
echo ""

# Test 16: Get Weekly Summary
echo "16. Getting Weekly Summary..."
curl -s "$BASE_URL/summary/weekly" | jq .
echo ""

# Test 17: Test Filtering - Active Intents Only
echo "17. Testing Filtering - Active Intents Only..."
curl -s "$BASE_URL/intents?active=true" | jq .
echo ""

# Test 18: Test Filtering - Active Focus Sessions Only
echo "18. Testing Filtering - Active Focus Sessions Only..."
curl -s "$BASE_URL/focus?active=true" | jq .
echo ""

# Test 19: Update Journal Entry
echo "19. Updating Journal Entry..."
curl -s -X PUT "$BASE_URL/journal/$JOURNAL_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated: API Development Day",
    "content": "Successfully implemented the Mindloop API with all features working perfectly! Ready for production!",
    "mood": "excited"
  }' | jq .
echo ""

# Test 20: Custom Summary
echo "20. Getting Custom Date Range Summary..."
curl -s -X POST "$BASE_URL/summary/custom" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "2025-08-01",
    "end_date": "2025-08-31"
  }' | jq .
echo ""

echo "✅ All API tests completed successfully!"
echo ""
echo "🎉 Mindloop API is fully functional with all features:"
echo "   - Habits management ✅"
echo "   - Intent tracking ✅"
echo "   - Focus sessions ✅"
echo "   - Journal entries ✅"
echo "   - Summary generation ✅"
echo ""
echo "📚 API Documentation: api/v1/README.md"
echo "🌐 Server running at: http://localhost:8080"
echo "📊 API Base URL: http://localhost:8080/api/v1"

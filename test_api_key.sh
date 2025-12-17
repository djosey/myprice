#!/bin/zsh
# Test if Anthropic API key is valid

if [ -z "$ANTHROPIC_API_KEY" ]; then
  echo "❌ ANTHROPIC_API_KEY is not set"
  exit 1
fi

# Clean the key
API_KEY=$(echo "$ANTHROPIC_API_KEY" | tr -d ' "'\''')

echo "Testing API key..."
echo "Key starts with: ${API_KEY:0:10}..."
echo "Key length: ${#API_KEY}"
echo ""

# Test the API key with a simple request
RESPONSE=$(curl -s -X POST https://api.anthropic.com/v1/messages \
  -H "Content-Type: application/json" \
  -H "x-api-key: $API_KEY" \
  -H "anthropic-version: 2023-06-01" \
  -d '{
    "model": "claude-3-5-sonnet-20241022",
    "max_tokens": 10,
    "messages": [{"role": "user", "content": "Hi"}]
  }')

if echo "$RESPONSE" | grep -q "authentication_error"; then
  echo "❌ API key is INVALID"
  echo ""
  echo "Response:"
  echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
  echo ""
  echo "Possible issues:"
  echo "  - Key is expired or revoked"
  echo "  - Key has extra characters/spaces"
  echo "  - Wrong key copied"
  echo ""
  echo "Get a new key from: https://console.anthropic.com/"
elif echo "$RESPONSE" | grep -q "id"; then
  echo "✅ API key is VALID!"
  echo ""
  echo "The key works. If you're still getting 401 errors, check:"
  echo "  1. Server is reading the key correctly (check server startup logs)"
  echo "  2. Key doesn't have quotes in ~/.zshrc"
else
  echo "⚠️  Unexpected response:"
  echo "$RESPONSE"
fi


#!/bin/bash
# Test script to verify Claude API is working

echo "=== Testing Claude API Integration ==="
echo ""

# Check if API key is set
if [ -z "$ANTHROPIC_API_KEY" ]; then
    echo "❌ ANTHROPIC_API_KEY is NOT set"
    echo "   Set it with: export ANTHROPIC_API_KEY='sk-ant-...'"
    exit 1
else
    echo "✅ ANTHROPIC_API_KEY is set (${#ANTHROPIC_API_KEY} chars)"
fi

echo ""
echo "Making test API call..."
echo ""

# Make API call and check response
RESPONSE=$(curl -s -X POST http://localhost:8080/api/analyze \
  -H "Content-Type: application/json" \
  -d '{"image_path": "/Users/donaldjosey/developer/myprice/PXL_20251206_222017258.jpg"}')

# Check if cart_description is in response
if echo "$RESPONSE" | grep -q "cart_description"; then
    echo "✅ SUCCESS: Claude is working! Found 'cart_description' in response"
    echo ""
    echo "Sample of response:"
    echo "$RESPONSE" | python3 -c "import sys,json; d=json.load(sys.stdin); llm=d.get('llm_output',{}); print('cart_description:', llm.get('cart_description', 'NOT FOUND')); print('item_categories:', llm.get('item_categories', 'NOT FOUND'))" 2>/dev/null || echo "Could not parse JSON"
elif echo "$RESPONSE" | grep -q "Parsed from Textract OCR output"; then
    echo "❌ FAILED: Using regex parser (Claude not working)"
    echo ""
    echo "Check your server terminal for one of these messages:"
    echo "  - 'Claude API not configured' → API key not set in server"
    echo "  - 'LLM parsing failed' → Claude API call failed"
else
    echo "⚠️  Could not determine status. Full response:"
    echo "$RESPONSE" | python3 -m json.tool | head -20
fi


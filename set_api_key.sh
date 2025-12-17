#!/bin/zsh
# Script to set Anthropic API key

echo "Setting ANTHROPIC_API_KEY..."
echo ""
echo "Enter your API key (it should start with 'sk-ant-'):"
read -s API_KEY

if [[ -z "$API_KEY" ]]; then
  echo "❌ No API key provided"
  exit 1
fi

if [[ ! "$API_KEY" =~ ^sk-ant- ]]; then
  echo "⚠️  Warning: API key doesn't start with 'sk-ant-'"
  echo "   Are you sure this is correct?"
  read -q "REPLY?Continue anyway? (y/n) "
  echo ""
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 1
  fi
fi

# Set for current session
export ANTHROPIC_API_KEY="$API_KEY"

# Add to ~/.zshrc if not already there
if ! grep -q "ANTHROPIC_API_KEY" ~/.zshrc 2>/dev/null; then
  echo "" >> ~/.zshrc
  echo "# Anthropic API Key" >> ~/.zshrc
  echo "export ANTHROPIC_API_KEY=\"$API_KEY\"" >> ~/.zshrc
  echo "✅ Added to ~/.zshrc (permanent)"
else
  echo "⚠️  ANTHROPIC_API_KEY already exists in ~/.zshrc"
  echo "   You may want to update it manually"
fi

echo "✅ API key set for current session"
echo ""
echo "To verify:"
echo "  echo \$ANTHROPIC_API_KEY"
echo ""
echo "To restart server:"
echo "  pkill -f myprice-api && DISABLE_CACHE=true ./myprice-api"

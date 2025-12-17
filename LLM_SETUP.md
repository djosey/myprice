# LLM Receipt Parsing Setup

The receipt analysis now uses **Claude API** for intelligent parsing instead of regex.

## Setup

1. **Get your Anthropic API key:**
   - Sign up at https://console.anthropic.com/
   - Create an API key
   - Copy the key

2. **Set the environment variable:**
   ```bash
   export ANTHROPIC_API_KEY="sk-ant-..."
   ```

   Or add it to your shell profile (`~/.zshrc` or `~/.bashrc`):
   ```bash
   echo 'export ANTHROPIC_API_KEY="sk-ant-..."' >> ~/.zshrc
   source ~/.zshrc
   ```

3. **Restart the API server:**
   ```bash
   ./myprice-api
   ```

## How It Works

1. **Image Upload** → Saved to `uploads/`
2. **AWS Textract** → OCR text extraction
3. **Claude API** → Intelligent parsing with:
   - Image + OCR text as input
   - Structured JSON output
   - Error correction
   - Context understanding

## Fallback

If `ANTHROPIC_API_KEY` is not set, the system falls back to the regex parser (less accurate).

## Testing

```bash
# Make sure API key is set
echo $ANTHROPIC_API_KEY

# Test the API
curl -X POST http://localhost:8080/api/analyze \
  -H "Content-Type: application/json" \
  -d '{"image_path": "/path/to/receipt.jpg"}' \
  | python3 -m json.tool
```

## Cost

Claude Sonnet 4 costs approximately:
- **$0.003 per image** (input)
- **$0.015 per 1K tokens** (output)

A typical receipt analysis costs **~$0.01-0.03** per receipt.


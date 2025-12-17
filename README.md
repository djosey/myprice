# myprice

A Go MCP (Model Context Protocol) server for multimodal receipt processing.

## Overview

This project provides an MCP server that exposes tools for:
- Loading receipt images (base64-encoded with MIME type)
- Parsing AWS Textract OCR output
- Writing structured receipt data to disk

The server is designed to be used with an LLM that orchestrates receipt extraction by:
1. Loading an image
2. Loading Textract OCR results
3. Comparing and correcting OCR text against visual content
4. Outputting normalized, structured receipt data

## Project Structure

```
myprice/
├── main.go                    # MCP server entrypoint
├── go.mod                     # Go module definition
├── tools/
│   ├── load_image.go          # load_image tool implementation
│   ├── load_textract.go       # load_textract tool implementation
│   └── write_output.go        # write_output tool implementation
├── internal/
│   └── receipt/
│       ├── schema.go          # Receipt output schema
│       └── normalize.go       # Text normalization helpers
└── README.md
```

## Building

```bash
go build -o myprice-mcp .
```

## MCP Tools

### `load_image`

Load an image file and return its base64-encoded bytes along with MIME type.

**Input:**
```json
{
  "path": "/path/to/receipt.jpg"
}
```

**Output:**
- Image content for visual inspection
- Structured metadata: `{ base64_data, mime_type, file_path, size_bytes }`

### `load_textract`

Load and parse an AWS Textract JSON output file.

**Input:**
```json
{
  "path": "/path/to/textract_output.json"
}
```

**Output:**
```json
{
  "page_count": 1,
  "lines": [
    { "text": "STORE NAME", "confidence": 99.5, "top": 0.12, "left": 0.35 },
    { "text": "$12.99", "confidence": 98.2, "top": 0.45, "left": 0.60 }
  ],
  "total_lines": 42,
  "file_path": "/path/to/textract_output.json"
}
```

### `write_output`

Write structured JSON data to a file.

**Input:**
```json
{
  "path": "/path/to/output.json",
  "data": {
    "vendor": "Target",
    "date": "2024-12-06",
    "items": [
      { "name": "Milk", "qty": 1, "price": 4.99 }
    ],
    "subtotal": 4.99,
    "tax": 0.40,
    "total": 5.39
  }
}
```

**Output:**
```json
{
  "success": true,
  "file_path": "/path/to/output.json",
  "bytes_written": 256
}
```

## Receipt Output Schema

The expected structured output for receipts:

```json
{
  "vendor": "Store Name",
  "date": "YYYY-MM-DD",
  "items": [
    { "name": "Item Name", "qty": 1, "price": 0.00 }
  ],
  "subtotal": 0.00,
  "tax": 0.00,
  "total": 0.00,
  "confidence_notes": "Any notes about OCR quality or corrections made",
  "anomalies": ["List of detected issues or inconsistencies"]
}
```

## Configuration

To use this MCP server with Claude Desktop or other MCP clients, add to your MCP config:

```json
{
  "mcpServers": {
    "myprice": {
      "command": "/path/to/myprice-mcp",
      "args": []
    }
  }
}
```

## Running Textract

Use the included `detect.sh` script to run AWS Textract on a receipt image:

```bash
./detect.sh
```

This requires AWS CLI configured with appropriate credentials.

## Development

### Prerequisites
- Go 1.22+
- AWS CLI (for Textract)

### Dependencies
- `github.com/modelcontextprotocol/go-sdk` - MCP Go SDK

### Testing
```bash
go test ./...
```

## License

MIT

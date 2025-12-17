// Package tools provides MCP tool implementations for receipt processing.
package tools

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// LoadImageInput defines the input parameters for load_image tool.
type LoadImageInput struct {
	Path string `json:"path" doc:"Absolute or relative path to the image file"`
}

// LoadImageOutput defines the output structure for load_image tool.
type LoadImageOutput struct {
	Base64Data string `json:"base64_data"`
	MimeType   string `json:"mime_type"`
	FilePath   string `json:"file_path"`
	SizeBytes  int64  `json:"size_bytes"`
}

// LoadImageTool returns the MCP tool definition for load_image.
func LoadImageTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "load_image",
		Description: "Load an image file and return its base64-encoded bytes along with MIME type. Useful for visual inspection of receipts.",
	}
}

// HandleLoadImage processes the load_image tool call.
// The handler returns both a CallToolResult (with image content) and the structured output.
func HandleLoadImage(ctx context.Context, req *mcp.CallToolRequest, input LoadImageInput) (*mcp.CallToolResult, LoadImageOutput, error) {
	if input.Path == "" {
		return nil, LoadImageOutput{}, fmt.Errorf("path is required")
	}

	// Read the file
	data, err := os.ReadFile(input.Path)
	if err != nil {
		return nil, LoadImageOutput{}, fmt.Errorf("failed to read image: %w", err)
	}

	// Get file info for size
	info, err := os.Stat(input.Path)
	if err != nil {
		return nil, LoadImageOutput{}, fmt.Errorf("failed to stat file: %w", err)
	}

	// Determine MIME type from extension
	ext := strings.ToLower(filepath.Ext(input.Path))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		// Fallback for common image types
		switch ext {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".png":
			mimeType = "image/png"
		case ".gif":
			mimeType = "image/gif"
		case ".webp":
			mimeType = "image/webp"
		case ".heic", ".heif":
			mimeType = "image/heic"
		default:
			mimeType = "application/octet-stream"
		}
	}

	// Encode to base64
	base64Data := base64.StdEncoding.EncodeToString(data)

	output := LoadImageOutput{
		Base64Data: base64Data,
		MimeType:   mimeType,
		FilePath:   input.Path,
		SizeBytes:  info.Size(),
	}

	// Return the image as content for the LLM to see
	// Note: ImageContent.Data takes raw bytes; the SDK handles base64 encoding
	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.ImageContent{
				Data:     data, // raw bytes, SDK encodes to base64
				MIMEType: mimeType,
			},
		},
	}

	return result, output, nil
}

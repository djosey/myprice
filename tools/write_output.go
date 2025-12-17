// Package tools provides MCP tool implementations for receipt processing.
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// WriteOutputInput defines the input parameters for write_output tool.
type WriteOutputInput struct {
	Path string `json:"path" doc:"Path where the JSON output should be written"`
	Data any    `json:"data" doc:"The structured data to write as JSON"`
}

// WriteOutputOutput defines the result of a write operation.
type WriteOutputOutput struct {
	Success      bool   `json:"success"`
	FilePath     string `json:"file_path"`
	BytesWritten int    `json:"bytes_written"`
}

// WriteOutputTool returns the MCP tool definition for write_output.
func WriteOutputTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "write_output",
		Description: "Write structured JSON data to a file. Use this to save the final parsed receipt data or intermediate results.",
	}
}

// HandleWriteOutput processes the write_output tool call.
func HandleWriteOutput(ctx context.Context, req *mcp.CallToolRequest, input WriteOutputInput) (*mcp.CallToolResult, WriteOutputOutput, error) {
	if input.Path == "" {
		return nil, WriteOutputOutput{}, fmt.Errorf("path is required")
	}

	if input.Data == nil {
		return nil, WriteOutputOutput{}, fmt.Errorf("data is required")
	}

	// Ensure the directory exists
	dir := filepath.Dir(input.Path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, WriteOutputOutput{}, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Serialize the data with pretty printing
	jsonData, err := json.MarshalIndent(input.Data, "", "  ")
	if err != nil {
		return nil, WriteOutputOutput{}, fmt.Errorf("failed to serialize data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(input.Path, jsonData, 0644); err != nil {
		return nil, WriteOutputOutput{}, fmt.Errorf("failed to write file: %w", err)
	}

	output := WriteOutputOutput{
		Success:      true,
		FilePath:     input.Path,
		BytesWritten: len(jsonData),
	}

	return nil, output, nil
}

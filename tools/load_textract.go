// Package tools provides MCP tool implementations for receipt processing.
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TextractBlock represents a single block from AWS Textract output.
type TextractBlock struct {
	BlockType     string          `json:"BlockType"`
	Confidence    float64         `json:"Confidence,omitempty"`
	Text          string          `json:"Text,omitempty"`
	ID            string          `json:"Id"`
	Geometry      *BlockGeometry  `json:"Geometry,omitempty"`
	Relationships []Relationship  `json:"Relationships,omitempty"`
}

// BlockGeometry contains position information for a block.
type BlockGeometry struct {
	BoundingBox *BoundingBox `json:"BoundingBox,omitempty"`
}

// BoundingBox defines the rectangular area of a block.
type BoundingBox struct {
	Width  float64 `json:"Width"`
	Height float64 `json:"Height"`
	Left   float64 `json:"Left"`
	Top    float64 `json:"Top"`
}

// Relationship defines connections between blocks.
type Relationship struct {
	Type string   `json:"Type"`
	IDs  []string `json:"Ids"`
}

// TextractDocument represents the full Textract response.
type TextractDocument struct {
	DocumentMetadata struct {
		Pages int `json:"Pages"`
	} `json:"DocumentMetadata"`
	Blocks []TextractBlock `json:"Blocks"`
}

// TextractLine represents a line of text with confidence and position.
type TextractLine struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
	Top        float64 `json:"top"`
	Left       float64 `json:"left"`
}

// LoadTextractInput defines the input parameters for load_textract tool.
type LoadTextractInput struct {
	Path string `json:"path" doc:"Path to the Textract JSON output file"`
}

// LoadTextractOutput is the simplified output for the LLM.
type LoadTextractOutput struct {
	PageCount  int            `json:"page_count"`
	Lines      []TextractLine `json:"lines"`
	TotalLines int            `json:"total_lines"`
	FilePath   string         `json:"file_path"`
}

// LoadTextractTool returns the MCP tool definition for load_textract.
func LoadTextractTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "load_textract",
		Description: "Load and parse an AWS Textract JSON output file. Returns extracted text lines with confidence scores and positions, sorted by vertical position (top to bottom).",
	}
}

// HandleLoadTextract processes the load_textract tool call.
func HandleLoadTextract(ctx context.Context, req *mcp.CallToolRequest, input LoadTextractInput) (*mcp.CallToolResult, LoadTextractOutput, error) {
	if input.Path == "" {
		return nil, LoadTextractOutput{}, fmt.Errorf("path is required")
	}

	// Read the file
	data, err := os.ReadFile(input.Path)
	if err != nil {
		return nil, LoadTextractOutput{}, fmt.Errorf("failed to read Textract file: %w", err)
	}

	// Parse the Textract JSON
	var doc TextractDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, LoadTextractOutput{}, fmt.Errorf("failed to parse Textract JSON: %w", err)
	}

	// Extract LINE blocks
	lines := make([]TextractLine, 0)
	for _, block := range doc.Blocks {
		if block.BlockType == "LINE" && block.Text != "" {
			line := TextractLine{
				Text:       block.Text,
				Confidence: block.Confidence,
			}
			if block.Geometry != nil && block.Geometry.BoundingBox != nil {
				line.Top = block.Geometry.BoundingBox.Top
				line.Left = block.Geometry.BoundingBox.Left
			}
			lines = append(lines, line)
		}
	}

	// Sort lines by vertical position (top to bottom), then by left position
	sort.Slice(lines, func(i, j int) bool {
		if lines[i].Top != lines[j].Top {
			return lines[i].Top < lines[j].Top
		}
		return lines[i].Left < lines[j].Left
	})

	output := LoadTextractOutput{
		PageCount:  doc.DocumentMetadata.Pages,
		Lines:      lines,
		TotalLines: len(lines),
		FilePath:   input.Path,
	}

	return nil, output, nil
}

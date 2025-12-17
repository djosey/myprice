// Package server provides LLM integration for receipt parsing.
package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"myprice/tools"
)

// ClaudeAPI handles calls to Anthropic's Claude API.
type ClaudeAPI struct {
	apiKey string
	client *http.Client
}

// NewClaudeAPI creates a new Claude API client.
func NewClaudeAPI() (*ClaudeAPI, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}

	// Clean up the API key (remove quotes, whitespace, etc.)
	apiKey = strings.TrimSpace(apiKey)
	apiKey = strings.Trim(apiKey, `"'`)

	// Validate API key format
	if !strings.HasPrefix(apiKey, "sk-ant-") {
		return nil, fmt.Errorf("API key format invalid: must start with 'sk-ant-' (got: %s...)", apiKey[:min(10, len(apiKey))])
	}

	if len(apiKey) < 20 {
		return nil, fmt.Errorf("API key too short (length: %d)", len(apiKey))
	}

	log.Printf("Claude API key loaded: %s... (length: %d)", apiKey[:10], len(apiKey))

	return &ClaudeAPI{
		apiKey: apiKey,
		client: &http.Client{},
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ReceiptOutput represents the structured receipt output from the LLM.
type ReceiptOutput struct {
	Vendor          string   `json:"vendor"`
	VendorFull      string   `json:"vendor_full,omitempty"`
	Address         string   `json:"address,omitempty"`
	Date            string   `json:"date"`
	Time            string   `json:"time,omitempty"`
	Items           []Item   `json:"items"`
	Fees            []Fee    `json:"fees,omitempty"`
	Subtotal        float64  `json:"subtotal"`
	Tax             float64  `json:"tax"`
	Total           float64  `json:"total"`
	Server          string   `json:"server,omitempty"`
	CheckNumber     string   `json:"check_number,omitempty"`
	Table           string   `json:"table,omitempty"`
	Customer        string   `json:"customer,omitempty"`
	CartDescription string   `json:"cart_description,omitempty"`
	ItemCategories  []string `json:"item_categories,omitempty"`
	ConfidenceNotes string   `json:"confidence_notes"`
	Anomalies       []string `json:"anomalies"`
}

// Item represents a line item on the receipt.
type Item struct {
	Name  string  `json:"name"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
}

// Fee represents a fee or surcharge on the receipt.
type Fee struct {
	Name   string  `json:"name"`
	Rate   string  `json:"rate,omitempty"`
	Amount float64 `json:"amount"`
}

// ParseReceiptWithLLM uses Claude API to parse receipt from image and OCR text.
func (c *ClaudeAPI) ParseReceiptWithLLM(imagePath string, textractOutput tools.LoadTextractOutput) (*ReceiptOutput, error) {
	// Read and encode image
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)

	// Detect MIME type from file extension
	ext := filepath.Ext(imagePath)
	mediaType := mime.TypeByExtension(ext)
	if mediaType == "" {
		// Fallback to common image types
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg":
			mediaType = "image/jpeg"
		case ".png":
			mediaType = "image/png"
		case ".gif":
			mediaType = "image/gif"
		case ".webp":
			mediaType = "image/webp"
		default:
			mediaType = "image/jpeg" // Default fallback
		}
	}

	// Build OCR text summary
	ocrText := buildOCRText(textractOutput)

	// Build the prompt
	prompt := buildReceiptPrompt(ocrText)

	// Prepare Claude API request
	requestBody := map[string]interface{}{
		"model":      "claude-sonnet-4-20250514",
		"max_tokens": 4096,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "image",
						"source": map[string]interface{}{
							"type":       "base64",
							"media_type": mediaType,
							"data":       imageBase64,
						},
					},
					{
						"type": "text",
						"text": prompt,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make API call
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	log.Printf("Calling Claude API for receipt parsing...")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Claude API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResponse struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResponse.Content) == 0 {
		return nil, fmt.Errorf("empty response from Claude API")
	}

	// Extract JSON from response (may be wrapped in markdown code blocks or have extra text)
	jsonText := apiResponse.Content[0].Text
	jsonText = strings.TrimSpace(jsonText)

	// Remove markdown code blocks if present
	jsonText = strings.TrimPrefix(jsonText, "```json")
	jsonText = strings.TrimPrefix(jsonText, "```")
	jsonText = strings.TrimSuffix(jsonText, "```")
	jsonText = strings.TrimSpace(jsonText)

	// Try to find JSON object boundaries if there's extra text
	if strings.Contains(jsonText, "{") && strings.Contains(jsonText, "}") {
		start := strings.Index(jsonText, "{")
		end := strings.LastIndex(jsonText, "}") + 1
		if start >= 0 && end > start {
			jsonText = jsonText[start:end]
		}
	}

	// Parse JSON into ReceiptOutput
	var receipt ReceiptOutput
	if err := json.Unmarshal([]byte(jsonText), &receipt); err != nil {
		log.Printf("Failed to parse JSON response: %v", err)
		log.Printf("Response text: %s", jsonText)
		return nil, fmt.Errorf("failed to parse JSON from LLM response: %w", err)
	}

	log.Printf("Successfully parsed receipt: vendor=%s, items=%d, total=$%.2f",
		receipt.Vendor, len(receipt.Items), receipt.Total)

	return &receipt, nil
}

// buildOCRText formats the Textract output into a readable text summary.
func buildOCRText(textract tools.LoadTextractOutput) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("OCR Results (%d lines, %d pages):\n\n", len(textract.Lines), textract.PageCount))

	for i, line := range textract.Lines {
		sb.WriteString(fmt.Sprintf("%d. [%.1f%% confidence] %s\n", i+1, line.Confidence, line.Text))
	}

	return sb.String()
}

// buildReceiptPrompt creates the prompt for Claude to parse the receipt.
func buildReceiptPrompt(ocrText string) string {
	return `You are a receipt parsing expert. Analyze the receipt image and OCR text to extract structured data.

**OCR Text Data:**
` + ocrText + `

**Instructions:**
1. Extract vendor information:
   - Vendor name (short/common name)
   - Vendor full name (if different from short name)
   - Address (if present)

2. Extract date and time:
   - Date (normalize to ISO format: YYYY-MM-DD)
   - Time (if present, format as HH:MM AM/PM)

3. Extract all line items:
   - Item name (clean up OCR errors intelligently)
   - Quantity (if specified, default to 1)
   - Price (per item or total for that line)

4. Extract financial totals:
   - Subtotal
   - Tax
   - Fees (service fees, tips, surcharges, etc.)
   - Total

5. Extract context information (if present):
   - Server/waitstaff name
   - Table number
   - Check/receipt number
   - Customer name

6. Handle OCR errors intelligently:
   - Correct obvious typos (e.g., "T0AST" → "TOAST", "Patr0n" → "Patron")
   - Use context to disambiguate (e.g., "3 Patron Silver" likely means qty=3)
   - Match item names with prices even if they're on different lines
   - Handle multi-line item names

7. Note any anomalies or low-confidence extractions in the anomalies array.

8. Generate a cart description:
   - Write a brief narrative description (2-4 sentences) summarizing what was purchased
   - Describe the shopping pattern or theme (e.g., "Weekly grocery shopping with focus on fresh produce and dairy", "Quick convenience store stop for snacks and beverages", "Restaurant meal with multiple courses and drinks")
   - Include context about the type of purchase (grocery shopping, restaurant meal, convenience store, etc.)

9. Categorize the items:
   - Identify the main categories/types of items purchased
   - Use common categories like: produce, dairy, meat, seafood, beverages, snacks, frozen, bakery, deli, prepared_foods, alcohol, household, personal_care, etc.
   - Include all relevant categories (items can belong to multiple categories)
   - Return as an array of category strings

**Output Format (JSON only, no markdown):**
{
  "vendor": "string",
  "vendor_full": "string (optional)",
  "address": "string (optional)",
  "date": "YYYY-MM-DD",
  "time": "HH:MM AM/PM (optional)",
  "items": [
    {"name": "string", "qty": number, "price": number}
  ],
  "fees": [
    {"name": "string", "rate": "string (optional)", "amount": number}
  ],
  "subtotal": number,
  "tax": number,
  "total": number,
  "server": "string (optional)",
  "check_number": "string (optional)",
  "table": "string (optional)",
  "customer": "string (optional)",
  "cart_description": "string - brief narrative description of the shopping cart/purchase (2-4 sentences)",
  "item_categories": ["string array of item categories like: produce, dairy, meat, beverages, snacks, etc."],
  "confidence_notes": "string describing confidence level and any issues",
  "anomalies": ["string array of any anomalies or uncertainties"]
}

**CRITICAL:** Return ONLY valid JSON. Do not include markdown code blocks, explanations, or any text before or after the JSON. Start with { and end with }.`
}

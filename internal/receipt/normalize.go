// Package receipt provides normalization helpers for receipt data.
package receipt

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	// pricePattern matches common price formats like $1.99, 1.99, $1,234.56
	pricePattern = regexp.MustCompile(`^\$?[\d,]+\.?\d*$`)

	// datePatterns for common receipt date formats
	datePatterns = []*regexp.Regexp{
		regexp.MustCompile(`\d{1,2}/\d{1,2}/\d{2,4}`),          // MM/DD/YYYY or M/D/YY
		regexp.MustCompile(`\d{1,2}-\d{1,2}-\d{2,4}`),          // MM-DD-YYYY
		regexp.MustCompile(`\d{4}-\d{2}-\d{2}`),                // YYYY-MM-DD
		regexp.MustCompile(`\w{3}\s+\d{1,2},?\s+\d{4}`),        // Jan 15, 2024
		regexp.MustCompile(`\d{1,2}\s+\w{3}\s+\d{4}`),          // 15 Jan 2024
	}
)

// NormalizePrice cleans a price string and parses it as a float64.
// Returns 0.0 if the string cannot be parsed.
func NormalizePrice(s string) float64 {
	// Remove dollar sign, commas, and whitespace
	cleaned := strings.TrimSpace(s)
	cleaned = strings.TrimPrefix(cleaned, "$")
	cleaned = strings.ReplaceAll(cleaned, ",", "")

	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0.0
	}
	return val
}

// IsPrice checks if a string looks like a price value.
func IsPrice(s string) bool {
	cleaned := strings.TrimSpace(s)
	return pricePattern.MatchString(cleaned)
}

// ExtractDate attempts to find a date string from text.
// Returns the matched date string or empty string if not found.
func ExtractDate(text string) string {
	for _, pattern := range datePatterns {
		if match := pattern.FindString(text); match != "" {
			return match
		}
	}
	return ""
}

// NormalizeVendorName cleans up a vendor name string.
func NormalizeVendorName(s string) string {
	// Trim whitespace and normalize to title case for consistency
	cleaned := strings.TrimSpace(s)
	return cleaned
}

// NormalizeItemName cleans up an item name.
func NormalizeItemName(s string) string {
	cleaned := strings.TrimSpace(s)
	// Remove common receipt artifacts
	cleaned = strings.TrimPrefix(cleaned, "*")
	cleaned = strings.TrimSuffix(cleaned, "*")
	return cleaned
}

// ParseQuantity attempts to extract a quantity from a string.
// Returns 1 as default if no quantity can be parsed.
func ParseQuantity(s string) int {
	cleaned := strings.TrimSpace(s)
	qty, err := strconv.Atoi(cleaned)
	if err != nil || qty < 1 {
		return 1
	}
	return qty
}




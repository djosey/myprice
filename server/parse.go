// Package server provides HTTP API endpoints for the receipt analysis tools.
package server

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	// Price patterns like $12.99, 12.99, $1,234.56
	priceRegex = regexp.MustCompile(`\$?([\d,]+\.?\d*)`)

	// Date patterns
	dateRegex = regexp.MustCompile(`\d{1,2}/\d{1,2}/\d{2,4}|\d{4}-\d{2}-\d{2}`)
)

// containsPrice checks if a string contains a price-like pattern.
func containsPrice(s string) bool {
	return strings.Contains(s, "$") || priceRegex.MatchString(s)
}

// containsDate checks if a string contains a date pattern.
func containsDate(s string) bool {
	return dateRegex.MatchString(s)
}

// extractPrice extracts a numeric price from a string.
func extractPrice(s string) float64 {
	matches := priceRegex.FindStringSubmatch(s)
	if len(matches) < 2 {
		return 0
	}

	// Remove commas and parse
	priceStr := strings.ReplaceAll(matches[1], ",", "")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0
	}
	return price
}

// extractItemName extracts the item name from a line (removes the price part).
func extractItemName(s string) string {
	// Remove price portion
	name := priceRegex.ReplaceAllString(s, "")
	// Remove $ signs
	name = strings.ReplaceAll(name, "$", "")
	// Trim whitespace
	name = strings.TrimSpace(name)

	// Skip if it's just a number or too short
	if len(name) < 2 {
		return ""
	}

	return name
}


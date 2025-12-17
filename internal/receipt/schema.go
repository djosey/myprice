// Package receipt defines the structured output schema for parsed receipts.
package receipt

// Item represents a single line item on a receipt.
type Item struct {
	Name  string  `json:"name"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
}

// Receipt represents the normalized, structured output from receipt analysis.
type Receipt struct {
	Vendor          string   `json:"vendor"`
	Date            string   `json:"date"`
	Items           []Item   `json:"items"`
	Subtotal        float64  `json:"subtotal"`
	Tax             float64  `json:"tax"`
	Total           float64  `json:"total"`
	ConfidenceNotes string   `json:"confidence_notes"`
	Anomalies       []string `json:"anomalies"`
}

// NewReceipt creates a new Receipt with initialized slices.
func NewReceipt() *Receipt {
	return &Receipt{
		Items:     make([]Item, 0),
		Anomalies: make([]string, 0),
	}
}




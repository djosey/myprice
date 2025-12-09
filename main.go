package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Hardcoded CPI series IDs (exactly as provided)
var seriesIDs = []string{
	"APU0000702111",
	"APU0000703112",
	"APU0000708112",
	"APU0000709112",
	"APU0000710212",
	"APU0000711111",
	"APU0000711211",
	"APU0000712112",
	"APU0000714221",
	"APU0000715211",
	"APU0000717311",
	"APU0000FF1101",
	"APU0000FJ1101",
	"APU0000FL2101",
	"APU0000FS1101",
}

const blsBaseURL = "https://api.bls.gov/publicAPI/v2/timeseries/data/"

// Structures for decoding BLS "latest" response

type blsResponse struct {
	Status       string          `json:"status"`
	ResponseTime int             `json:"responseTime"`
	Message      []string        `json:"message"`
	Results      blsResultHolder `json:"Results"`
}

type blsResultHolder struct {
	Series []blsSeries `json:"series"`
}

type blsSeries struct {
	SeriesID string       `json:"seriesID"`
	Data     []blsDataRow `json:"data"`
}

type blsDataRow struct {
	Year       string       `json:"year"`
	Period     string       `json:"period"`
	PeriodName string       `json:"periodName"`
	Latest     string       `json:"latest,omitempty"`
	Value      string       `json:"value"`
	Footnotes  []blsFootnote `json:"footnotes"`
}

type blsFootnote struct {
	Code string `json:"code,omitempty"`
	Text string `json:"text,omitempty"`
}

// Our own simplified output type
type latestValue struct {
	SeriesID   string `json:"series_id"`
	Year       string `json:"year"`
	Period     string `json:"period"`
	PeriodName string `json:"period_name"`
	Value      string `json:"value"`
}

// HTTP client with timeout
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func main() {
	http.HandleFunc("/latest", latestHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting BLS latest-value service on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// latestHandler fetches the latest value for each hardcoded series ID
func latestHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("BLS_API_KEY")

	results := make([]latestValue, 0, len(seriesIDs))
	for _, id := range seriesIDs {
		lv, err := fetchLatestForSeries(id, apiKey)
		if err != nil {
			log.Printf("error fetching series %s: %v", id, err)
			// You can choose to skip or return an error; here we skip and continue.
			continue
		}
		results = append(results, lv)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// fetchLatestForSeries calls the BLS "latest" endpoint for a single series ID
func fetchLatestForSeries(seriesID, apiKey string) (latestValue, error) {
	// Build URL: /timeseries/data/{seriesID}?latest=true[&registrationkey=...]
	url := fmt.Sprintf("%s%s?latest=true", blsBaseURL, seriesID)
	if apiKey != "" {
		url = url + "&registrationkey=" + apiKey
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return latestValue{}, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", "bls-latest-service/1.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return latestValue{}, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return latestValue{}, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var br blsResponse
	if err := json.NewDecoder(resp.Body).Decode(&br); err != nil {
		return latestValue{}, fmt.Errorf("decode response: %w", err)
	}

	if br.Status != "REQUEST_SUCCEEDED" {
		return latestValue{}, fmt.Errorf("bls status %q, messages: %v", br.Status, br.Message)
	}

	if len(br.Results.Series) == 0 || len(br.Results.Series[0].Data) == 0 {
		return latestValue{}, fmt.Errorf("no data for series %s", seriesID)
	}

	// With latest=true, BLS returns only the latest data row in Results.Series[0].Data[0]
	d := br.Results.Series[0].Data[0]

	return latestValue{
		SeriesID:   seriesID,
		Year:       d.Year,
		Period:     d.Period,
		PeriodName: d.PeriodName,
		Value:      d.Value,
	}, nil
}

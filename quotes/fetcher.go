package quotes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Quote represents a quote from the API
type Quote struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

// Fetcher handles quote retrieval
type Fetcher struct {
	client  *http.Client
	baseURL string
}

// NewFetcher creates a new quote fetcher
func NewFetcher() *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.quotable.io",
	}
}

// FetchRandomQuote fetches a random quote from the API
func (f *Fetcher) FetchRandomQuote() (*Quote, error) {
	url := f.baseURL + "/random"

	resp, err := f.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quote: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var quote Quote
	if err := json.Unmarshal(body, &quote); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate quote content
	if quote.Content == "" {
		return nil, fmt.Errorf("received empty quote content")
	}

	return &quote, nil
}

// FetchRandomQuoteWithFallback fetches a quote with fallback to hardcoded quotes
func (f *Fetcher) FetchRandomQuoteWithFallback() *Quote {
	quote, err := f.FetchRandomQuote()
	if err != nil {
		// Return a fallback quote if API fails
		return &Quote{
			Content: "The quick brown fox jumps over the lazy dog.",
			Author:  "Fallback",
		}
	}
	return quote
}

// GetFallbackQuotes returns a list of hardcoded quotes for offline use
func GetFallbackQuotes() []Quote {
	return []Quote{
		{
			Content: "The quick brown fox jumps over the lazy dog.",
			Author:  "Typing Test",
		},
		{
			Content: "To be or not to be, that is the question.",
			Author:  "William Shakespeare",
		},
		{
			Content: "The only way to do great work is to love what you do.",
			Author:  "Steve Jobs",
		},
		{
			Content: "In the middle of difficulty lies opportunity.",
			Author:  "Albert Einstein",
		},
		{
			Content: "Success is not final, failure is not fatal: it is the courage to continue that counts.",
			Author:  "Winston Churchill",
		},
	}
}


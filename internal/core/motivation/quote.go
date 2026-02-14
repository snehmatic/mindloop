package motivation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Quote struct {
	Text       string `json:"q"`
	Author     string `json:"a"`
	HTML       string `json:"h"`
	RetryAfter int    `json:"retry_after,omitempty"`
}

// Configurable constants for rate limiting
var (
	// RateLimit is the maximum number of requests allowed within the RateWindow
	RateLimit = 5
	// RateWindow is the duration for the rate limit window
	RateWindow = 30 * time.Second
)

var (
	cacheMu           sync.Mutex
	cachedQuote       *Quote
	requestTimestamps []time.Time
)

// FetchRandomQuote fetches a random quote from the ZenQuotes API
// It implements a sliding window rate limiter.
func FetchRandomQuote() (*Quote, error) {
	cacheMu.Lock()
	now := time.Now()

	// Clean up old timestamps
	validTimestamps := make([]time.Time, 0, len(requestTimestamps))
	for _, t := range requestTimestamps {
		if now.Sub(t) <= RateWindow {
			validTimestamps = append(validTimestamps, t)
		}
	}
	requestTimestamps = validTimestamps

	// Check rate limit
	if len(requestTimestamps) >= RateLimit {
		// Rate limit reached
		// Calculate wait time based on the oldest timestamp in the window
		var waitTime int
		if len(requestTimestamps) > 0 {
			earliest := requestTimestamps[0]
			resetTime := earliest.Add(RateWindow)
			waitTime = int(time.Until(resetTime).Seconds())
			if waitTime < 1 {
				waitTime = 1
			}
		} else {
			waitTime = int(RateWindow.Seconds())
		}

		if cachedQuote != nil {
			// Return cached quote with RetryAfter
			// Create a copy to avoid mutating the cached pointer
			q := *cachedQuote
			q.RetryAfter = waitTime
			cacheMu.Unlock()
			return &q, nil
		}

		cacheMu.Unlock()
		return nil, fmt.Errorf("rate limit reached and no cache available")
	}
	cacheMu.Unlock()

	// Rate limit allowed, fetch fresh quote
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("https://zenquotes.io/api/random")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quote: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode == http.StatusTooManyRequests {
		// Handle external 429
		cacheMu.Lock()
		defer cacheMu.Unlock()
		if cachedQuote != nil {
			q := *cachedQuote
			q.RetryAfter = int(RateWindow.Seconds())
			return &q, nil
		}
		return nil, fmt.Errorf("upstream API rate limit reached")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var quotes []Quote
	if err := json.NewDecoder(resp.Body).Decode(&quotes); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(quotes) == 0 {
		return nil, fmt.Errorf("no quotes returned")
	}

	// Update cache and timestamps
	cacheMu.Lock()
	cachedQuote = &quotes[0]
	// Clear the RetryAfter from the cached version just in case (though it comes empty from API)
	cachedQuote.RetryAfter = 0
	requestTimestamps = append(requestTimestamps, time.Now())
	cacheMu.Unlock()

	return &quotes[0], nil
}

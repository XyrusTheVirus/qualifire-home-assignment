package transports

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"qualifire-home-assignment/internal/loggers"
	"qualifire-home-assignment/internal/models"
	"time"
)

// LoggingTransport is an http.RoundTripper that logs the request method and duration
type LoggingTransport struct {
	Base http.RoundTripper
	Req  models.ProxyRequest
}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Save the request before it consumed by the LLM provider client
	bodyBytes, _ := io.ReadAll(req.Body)
	// Create a map to hold the JSON data
	var data map[string]interface{}
	// Unmarshal the JSON into the map
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	start := time.Now()
	resp, err := t.base().RoundTrip(req)
	duration := time.Since(start)
	entry := &models.LogEntry{
		Timestamp:   time.Now().Format(time.RFC3339),
		Method:      req.Method,
		DurationMS:  duration.Milliseconds(),
		Provider:    t.Req.Provider,
		VirtualKey:  t.Req.VirtualKey,
		RequestBody: data,
	}

	l := loggers.Logger{Entry: entry}
	if err != nil {
		l.Entry.Status = 0
		l.Entry.Error = "transport error: " + err.Error()
		l.Error()
		return nil, err
	}

	l.Entry.Status = resp.StatusCode
	// Use TeeReader to capture streamed data as it’s read
	var buf bytes.Buffer
	resp.Body = io.NopCloser(io.TeeReader(resp.Body, &buf))

	// When the caller finishes reading, log asynchronously
	go func() {
		// Wait until caller closes Body (i.e., finishes reading)
		defer func() {
			res := truncate(buf.String(), 2048)
			if err = json.Unmarshal(res, &data); err != nil {
				return
			}
			l.Entry.ResponseBody = data
			if l.Entry.Status >= 400 {
				l.Error()
			} else {
				l.Info()
			}
		}()
	}()

	return resp, nil
}

// Fallback base transport
func (t *LoggingTransport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

// Optional truncation helper to avoid huge logs
func truncate(s string, limit int) []byte {
	if len(s) <= limit {
		return []byte(s)
	}
	return []byte(s[:limit] + "...(truncated)")
}

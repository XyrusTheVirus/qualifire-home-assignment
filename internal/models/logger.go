package models

type LogEntry struct {
	Timestamp    string                 `json:"timestamp"`
	VirtualKey   string                 `json:"virtual_key"`
	Provider     string                 `json:"provider"`
	Method       string                 `json:"method"`
	Status       int                    `json:"status"`
	DurationMS   int64                  `json:"duration_ms"`
	RequestBody  map[string]interface{} `json:"request,omitempty"`
	ResponseBody map[string]interface{} `json:"response,omitempty"`
	Error        string                 `json:"error,omitempty"`
}

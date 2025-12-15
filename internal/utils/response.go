package utils

// WebResponse defines the structure for standard web responses
// used for successful API responses with data.
type WebResponse struct {
	Code   int32       `json:"code"`   // HTTP status code
	Status string      `json:"status"` // Textual status message
	Data   interface{} `json:"data"`   // Payload data
}

// ErrorResponse defines the structure for error responses
// used when an API request fails.
type ErrorResponse struct {
	Code    int32       `json:"code"`    // HTTP status code
	Status  string      `json:"status"`  // Textual status message
	Message interface{} `json:"message"` // Error message or details
}

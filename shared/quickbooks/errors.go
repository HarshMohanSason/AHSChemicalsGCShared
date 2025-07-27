package quickbooks

type QBErrorResponse struct {
	Fault QBFault `json:"Fault"`
	Time  string  `json:"time,omitempty"`
}

type QBFault struct {
	ErrorList []QBError `json:"Error"`
	Type      string    `json:"type"`
}

type QBError struct {
	Message       string `json:"Message"`              // Human-readable error message
	Detail        string `json:"Detail,omitempty"`     // More details if available
	Code          string `json:"code"`                 // e.g., "6140", "6240", etc.
	Element       string `json:"element,omitempty"`    // Field or object that caused the error
	Severity      string `json:"severity"`             // e.g., "Error", "Warning"
	InnerErrorMsg string `json:"InnerError,omitempty"` // Rare - sometimes used in deeper errors
}
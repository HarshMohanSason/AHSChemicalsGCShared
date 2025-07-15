package quickbooks

//Represents the response from Quickbooks when an error occurs
type QBErrorResponse struct {
	Fault QBErrorFault `json:"Fault"`
	Time  string       `json:"time"`
}

type QBErrorFault struct {
	Error []QBErrorDetail `json:"Error"`
	Type  string          `json:"type"`
}

type QBErrorDetail struct {
	Message string `json:"Message"`
	Detail  string `json:"Detail"`
	Code    string `json:"code"`
	Element string `json:"element,omitempty"`
}
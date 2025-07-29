package qbmodels

type QBErrorResponse struct {
	Warnings           any     `json:"warnings"`
	IntuitObject       any     `json:"intuitObject"`
	Fault              QBFault `json:"fault"`
	Report             any     `json:"report"`
	QueryResponse      any     `json:"queryResponse"`
	BatchItemResponse  any     `json:"batchItemResponse"`
	AttachableResponse any     `json:"attachableResponse"`
	SyncErrorResponse  any     `json:"syncErrorResponse"`
	RequestId          string  `json:"requestId"`
	Time               int64   `json:"time"`
	Status             string  `json:"status"`
	CDCResponse        any     `json:"cdcresponse"`
}

type QBFault struct {
	Error []QBError `json:"error"`
}

type QBError struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Code    string `json:"code"`
	Element string `json:"element"`
}
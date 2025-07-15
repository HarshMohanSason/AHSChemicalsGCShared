package shared

import (
	"sync"
)

var (
	InitFirebaseOnce   sync.Once //Firebase initialize once
	InitQuickBooksOnce sync.Once //Quickbooks initialize once
	InitGCPOnce        sync.Once //Google cloud manager initialize once
	InitSendGridOnce   sync.Once //Sendgrid initialize once
	InitCompanyDetails sync.Once //Company details initialize once
)

// Meta data for json returned by quickbooks for each product/customer/invoice etc...
type MetaData struct {
	CreatedAt string `json:"CreateTime,omitempty"`
	UpdatedAt string `json:"LastUpdatedTime,omitempty"`
}

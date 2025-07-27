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
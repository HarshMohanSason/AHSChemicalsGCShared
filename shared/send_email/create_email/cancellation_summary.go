package create_email

import (
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/company_details"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/send_email"
)

func CreateCancellationSummaryEmail(time time.Time) *send_email.EmailMetaData{
	emailData := &send_email.EmailMetaData{
		Recipients: company_details.EMAILINTERNALRECIPENTS,
		Data: map[string]any{
			"month_year":  time.Format("January 2006"),
		},
		TemplateID:  send_email.CANCELLED_ORDER_SUMMARY_TEMPLATE_ID,
	}
	return emailData
}
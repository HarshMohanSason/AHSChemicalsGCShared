package create_email

import (
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/company_details"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/send_email"
)

func CreateContactUsAdminEmail(c *models.ContactUsForm, attachments []send_email.Attachment) *send_email.EmailMetaData {
	emailData := &send_email.EmailMetaData{
		Recipients: company_details.EMAILINTERNALRECIPENTS,
		Data: map[string]any{
			"name":     c.Name,
			"email":    c.Email,
			"phone":    c.Location,
			"location": c.Location,
			"message":  c.Message,
			"year":     time.Now().Year(),
		},
		TemplateID:  send_email.CONTACT_US_ADMIN_TEMPLATE_ID,
		Attachments: attachments,
	}
	return emailData
}

func CreateContactUsUserEmail(c *models.ContactUsForm, attachments []send_email.Attachment) *send_email.EmailMetaData {
	emailData := &send_email.EmailMetaData{
		Recipients: map[string]string{c.Email: c.Name},
		Data: map[string]any{
			"name":     c.Name,
			"email":    c.Email,
			"phone":    c.Location,
			"location": c.Location,
			"message":  c.Message,
		},
		TemplateID:  send_email.CONTACT_US_USER_TEMPLATE_ID,
		Attachments: attachments,
	}
	return emailData
}
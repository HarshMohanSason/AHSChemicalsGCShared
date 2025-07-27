package send_email

import (
	"errors"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/company_details"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendMail(e *EmailMetaData) (*rest.Response, error) {
	from := mail.NewEmail("AHSChemicals", company_details.COMPANYEMAIL)

	var recipients []*mail.Email
	for email, name := range e.Recipients {
		recipients = append(recipients, mail.NewEmail(name, email))
	}

	if len(recipients) == 0 {
		return nil, errors.New("No recipients found to send a mail")
	}

	if e.TemplateID == "" {
		return nil, errors.New("No template ID found for sending the mail")
	}

	// Personalization for dynamic template data.
	p := mail.NewPersonalization()
	p.AddTos(recipients...)

	// Key-value pairs used in the send grid dynamic template
	for key, value := range e.Data {
		p.SetDynamicTemplateData(key, value)
	}

	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.AddPersonalizations(p)
	message.SetTemplateID(e.TemplateID)

	// Add attachments if any
	for _, item := range e.Attachments {
		attachment := mail.NewAttachment()
		attachment.SetType(item.MimeType)
		attachment.SetContent(item.Base64Content)
		attachment.SetFilename(item.FileName)
		attachment.SetDisposition("attachment")
		message.AddAttachment(attachment)
	}

	client := sendgrid.NewSendClient(SENDGRID_API_KEY)
	response, err := client.Send(message)
	if err != nil {
		return nil, err
	}

	return response, nil
}
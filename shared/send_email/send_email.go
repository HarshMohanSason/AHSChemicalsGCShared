package send_email

import (
	"errors"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/company_details"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// EmailMetaData defines the structure of the email payload expected in the request body.
type EmailMetaData struct {
	Recipients map[string]string `json:"recipients"`  
	Data       map[string]any    `json:"data"` //Data contains the key value pairs used to send values to the template in sendgrid       
	TemplateID string            `json:"template_id"` // ID of the dynamic template in sendgrid used to send mails
}

func SendEmail(metaData EmailMetaData) error {

	from := mail.NewEmail("AHSChemicals", company_details.COMPANYEMAIL)

	var recipients []*mail.Email
	for email, name := range metaData.Recipients {
		recipients = append(recipients, mail.NewEmail(name, email))
	}

	if len(recipients) == 0 {
		return errors.New("No recipents found to send a mail")
	}

	if metaData.TemplateID == "" {
		return errors.New("No template ID found for sending the mail")
	}

	// Configure personalization for dynamic template data.
	p := mail.NewPersonalization()
	p.AddTos(recipients...)

	//Set the key and values for the dynamic template data
	for key, value := range metaData.Data {
		p.SetDynamicTemplateData(key, value)
	}

	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.AddPersonalizations(p)
	message.SetTemplateID(metaData.TemplateID)

	client := sendgrid.NewSendClient(SENDGRID_API_KEY)
	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}

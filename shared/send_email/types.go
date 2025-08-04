package send_email

type Attachment struct {
	Base64Content string `json:"name"`
	MimeType      string `json:"data"`
	FileName      string `json:"content_type"` //FileName of the attachment. Must have an extension because on some mail applications, the attachment cannot be viewed if it does not have an extension. Eg: "purchase_order.pdf"
}

type EmailMetaData struct {
	Recipients  map[string]string `json:"recipients"`  //Map of recipent emails as keys and names as values
	Data        map[string]any    `json:"data"`        //Key value pairs of data that will be used in the send grid dynamic template
	TemplateID  string            `json:"template_id"` //Sendgrid dynamic template id
	Attachments []Attachment      `json:"attachments"` //List of attachments for the email
}

func (m *EmailMetaData) AddAttachment(attachment Attachment) {
	m.Attachments = append(m.Attachments, attachment)
}
func (m *EmailMetaData) AddAttachments(attachments []Attachment) {
	m.Attachments = append(m.Attachments, attachments...)
}
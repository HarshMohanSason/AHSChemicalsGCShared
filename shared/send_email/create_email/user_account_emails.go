package create_email

import (
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/send_email"
)

func CreateUserAccountCreatedEmail(createdUser *models.UserAccountCreate) *send_email.EmailMetaData {
	emailData := &send_email.EmailMetaData{
		Recipients: map[string]string{createdUser.Email: createdUser.Name},
		Data: map[string]any{
			"name":     createdUser.Name,
			"email":    createdUser.Email,
			"password": createdUser.Password,
		},
		TemplateID:  send_email.ACCOUNT_CREATED_USER_TEMPLATE_ID,
	}
	return emailData
}

func CreateDeleteUserAccountEmail(email, name string) *send_email.EmailMetaData {
	emailData := &send_email.EmailMetaData{
		Recipients: map[string]string{email: name},
		Data: map[string]any{
			"name":  name,
			"email": email,
		},
		TemplateID:  send_email.ACCOUNT_DELETED_USER_TEMPLATE_ID,
	}
	return emailData
}

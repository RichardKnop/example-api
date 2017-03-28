package accounts

import (
	"bytes"
	"fmt"
	htmlTemplate "html/template"
	"io/ioutil"
	"strings"
	textTemplate "text/template"

	"github.com/RichardKnop/example-api/config"
	"github.com/RichardKnop/example-api/services/email"
	"github.com/RichardKnop/example-api/models"
)

var (
	htmlEmailLayout                       = "./accounts/templates/email_layout.html"
	htmlEmailStyles                       = "./accounts/templates/styles.css"
	confirmEmailTemplateHTML              = "./accounts/templates/confirm_email.html"
	confirmEmailTemplateTxt               = "./accounts/templates/confirm_email.txt"
	passwordResetEmailTemplateHTML        = "./accounts/templates/password_reset_email.html"
	passwordResetEmailTemplateTxt         = "./accounts/templates/password_reset_email.txt"
	invitationEmailTemplateHTML           = "./accounts/templates/invitation_email.html"
	invitationEmailTemplateTxt            = "./accounts/templates/invitation_email.txt"
	onboardingCheckpointEmailTemplateHTML = "./accounts/templates/onboarding_checkpoint_email.html"
	onboardingCheckpointEmailTemplateTxt  = "./accounts/templates/onboarding_checkpoint_email.txt"
)

// EmailFactory facilitates construction of email.Email objects
type EmailFactory struct {
	cnf *config.Config
}

// NewEmailFactory starts a new emailFactory instance
func NewEmailFactory(cnf *config.Config) *EmailFactory {
	return &EmailFactory{cnf: cnf}
}

// NewConfirmationEmail returns a confirmation email
func (f *EmailFactory) NewConfirmationEmail(o *models.Confirmation) (*email.Message, error) {
	// Define a greetings name for the user
	name := o.User.GetName()
	if name == "" {
		name = "there"
	}
	name = strings.Split(name, " ")[0]

	// Confirmation link where the user can confirm his/her email
	link := fmt.Sprintf(
		f.cnf.AppSpecific.ConfirmationURLFormat,
		f.cnf.Web.AppScheme,
		f.cnf.Web.AppHost,
		o.Reference,
	)

	// The email subject
	subject := "Please confirm your email address"

	// Plain text email
	plainTextContent, err := newConfirmationEmailPlainTextContent(
		subject,
		name,
		f.cnf.AppSpecific.CompanyName,
		link,
	)
	if err != nil {
		return nil, err
	}

	// Read CSS styles file
	inlineStyles, err := ioutil.ReadFile(htmlEmailStyles)
	if err != nil {
		return nil, err
	}

	// HTML email
	htmlContent, err := newConfirmationEmailHTMLContent(
		subject,
		string(inlineStyles),
		name,
		f.cnf.AppSpecific.CompanyName,
		link,
	)
	if err != nil {
		return nil, err
	}

	return &email.Message{
		Subject: subject,
		Recipients: []*email.Recipient{&email.Recipient{
			Email: email.Email{
				Address: o.User.OauthUser.Username,
				Name:    o.User.GetName(),
			},
		}},
		From: &email.Sender{
			Email: email.Email{
				Address: f.cnf.AppSpecific.CompanyNoreplyEmail,
				Name:    f.cnf.AppSpecific.CompanyName,
			},
		},
		Text: plainTextContent,
		HTML: htmlContent,
	}, nil
}

func newConfirmationEmailPlainTextContent(title, name, company, link string) (string, error) {
	templateContents, err := ioutil.ReadFile(confirmEmailTemplateTxt)
	if err != nil {
		return "", err
	}
	tmpl, err := textTemplate.New(confirmEmailTemplateTxt).Parse(string(templateContents))
	if err != nil {
		return "", err
	}

	inventory := map[string]interface{}{
		"title":      title,
		"name":       name,
		"company":    company,
		"confirmURL": link,
	}

	var parsedTemplate bytes.Buffer
	if err := tmpl.Execute(&parsedTemplate, inventory); err != nil {
		return "", err
	}

	return parsedTemplate.String(), nil
}

func newConfirmationEmailHTMLContent(title, inlineStyles, name, company, link string) (string, error) {
	// Layout
	layoutContents, err := ioutil.ReadFile(htmlEmailLayout)
	if err != nil {
		return "", err
	}
	layoutTmpl, err := htmlTemplate.New(htmlEmailLayout).Parse(string(layoutContents))
	if err != nil {
		return "", err
	}

	// Content
	templateContents, err := ioutil.ReadFile(confirmEmailTemplateHTML)
	if err != nil {
		return "", err
	}
	tmpl, err := htmlTemplate.New(confirmEmailTemplateHTML).Parse(string(templateContents))
	if err != nil {
		return "", err
	}

	var (
		inventory                   map[string]interface{}
		parsedContent, parsedLayout bytes.Buffer
	)

	// Parse the content template
	inventory = map[string]interface{}{
		"name":       name,
		"company":    company,
		"confirmURL": link,
	}
	if err := tmpl.Execute(&parsedContent, inventory); err != nil {
		return "", err
	}

	// Insert the content into the layout
	inventory = map[string]interface{}{
		"title":        title,
		"inlineStyles": htmlTemplate.CSS(inlineStyles),
		"content":      htmlTemplate.HTML(parsedContent.String()),
		"company":      company,
	}
	if err := layoutTmpl.Execute(&parsedLayout, inventory); err != nil {
		return "", err
	}

	return parsedLayout.String(), nil
}

// NewPasswordResetEmail returns a password reset email
func (f *EmailFactory) NewPasswordResetEmail(o *models.PasswordReset) (*email.Message, error) {
	// Define a greetings name for the user
	name := o.User.GetName()
	if name == "" {
		name = "friend"
	}
	name = strings.Split(name, " ")[0]

	// Password reset link where the user can set a new password
	link := fmt.Sprintf(
		f.cnf.AppSpecific.PasswordResetURLFormat,
		f.cnf.Web.AppScheme,
		f.cnf.Web.AppHost,
		o.Reference,
	)

	// The email subject
	subject := fmt.Sprintf("Reset password for %s", f.cnf.Web.AppHost)

	// Plain text email
	plainTextContent, err := newPasswordResetEmailPlainTextContent(
		subject,
		name,
		f.cnf.AppSpecific.CompanyName,
		link,
	)
	if err != nil {
		return nil, err
	}

	// Read CSS styles file
	inlineStyles, err := ioutil.ReadFile(htmlEmailStyles)
	if err != nil {
		return nil, err
	}

	// HTML email
	htmlContent, err := newPasswordResetEmailHTMLContent(
		subject,
		string(inlineStyles),
		name,
		f.cnf.AppSpecific.CompanyName,
		link,
	)
	if err != nil {
		return nil, err
	}

	return &email.Message{
		Subject: subject,
		Recipients: []*email.Recipient{&email.Recipient{
			Email: email.Email{
				Address: o.User.OauthUser.Username,
				Name:    o.User.GetName(),
			},
		}},
		From: &email.Sender{
			Email: email.Email{
				Address: f.cnf.AppSpecific.CompanyNoreplyEmail,
				Name:    f.cnf.AppSpecific.CompanyName,
			},
		},
		Text: plainTextContent,
		HTML: htmlContent,
	}, nil
}

func newPasswordResetEmailPlainTextContent(title, name, company, link string) (string, error) {
	templateContents, err := ioutil.ReadFile(passwordResetEmailTemplateTxt)
	if err != nil {
		return "", err
	}
	tmpl, err := textTemplate.New(confirmEmailTemplateTxt).Parse(string(templateContents))
	if err != nil {
		return "", err
	}

	inventory := map[string]interface{}{
		"title":            title,
		"name":             name,
		"company":          company,
		"passwordResetURL": link,
	}

	var parsedTemplate bytes.Buffer
	if err := tmpl.Execute(&parsedTemplate, inventory); err != nil {
		return "", err
	}

	return parsedTemplate.String(), nil
}

func newPasswordResetEmailHTMLContent(title, inlineStyles, name, company, link string) (string, error) {
	// Layout
	layoutContents, err := ioutil.ReadFile(htmlEmailLayout)
	if err != nil {
		return "", err
	}
	layoutTmpl, err := htmlTemplate.New(htmlEmailLayout).Parse(string(layoutContents))
	if err != nil {
		return "", err
	}

	// Content
	templateContents, err := ioutil.ReadFile(passwordResetEmailTemplateHTML)
	if err != nil {
		return "", err
	}
	tmpl, err := htmlTemplate.New(passwordResetEmailTemplateHTML).Parse(string(templateContents))
	if err != nil {
		return "", err
	}

	var (
		inventory                   map[string]interface{}
		parsedContent, parsedLayout bytes.Buffer
	)

	// Parse the content template
	inventory = map[string]interface{}{
		"name":             name,
		"company":          company,
		"passwordResetURL": link,
	}
	if err := tmpl.Execute(&parsedContent, inventory); err != nil {
		return "", err
	}

	// Insert the content into the layout
	inventory = map[string]interface{}{
		"title":        title,
		"inlineStyles": htmlTemplate.CSS(inlineStyles),
		"content":      htmlTemplate.HTML(parsedContent.String()),
		"company":      company,
	}
	if err := layoutTmpl.Execute(&parsedLayout, inventory); err != nil {
		return "", err
	}

	return parsedLayout.String(), nil
}

// NewInvitationEmail returns a user invite email
func (f *EmailFactory) NewInvitationEmail(o *models.Invitation) (*email.Message, error) {
	// Define a greetings name for the invited user
	name := o.InvitedUser.GetName()
	if name == "" {
		name = "friend"
	}
	name = strings.Split(name, " ")[0]

	// Define a name of the person who invited the new user
	invitedBy := o.InvitedByUser.GetName()
	if invitedBy == "" {
		invitedBy = o.InvitedByUser.OauthUser.Username
	}

	// Confirmation link where the invited user can set his/her password
	link := fmt.Sprintf(
		f.cnf.AppSpecific.InvitationURLFormat,
		f.cnf.Web.AppScheme,
		f.cnf.Web.AppHost,
		o.Reference,
	)

	// The email subject
	subject := fmt.Sprintf("You have been invited to join %s", f.cnf.Web.AppHost)

	// Plain text email
	plainTextContent, err := newInvitationEmailPlainTextContent(
		subject,
		name,
		f.cnf.AppSpecific.CompanyName,
		invitedBy,
		link,
	)
	if err != nil {
		return nil, err
	}

	// Read CSS styles file
	inlineStyles, err := ioutil.ReadFile(htmlEmailStyles)
	if err != nil {
		return nil, err
	}

	// HTML email
	htmlContent, err := newInvitationEmailHTMLContent(
		subject,
		string(inlineStyles),
		name,
		f.cnf.AppSpecific.CompanyName,
		invitedBy,
		link,
	)
	if err != nil {
		return nil, err
	}

	return &email.Message{
		Subject: subject,
		Recipients: []*email.Recipient{&email.Recipient{
			Email: email.Email{
				Address: o.InvitedUser.OauthUser.Username,
				Name:    o.InvitedUser.GetName(),
			},
		}},
		From: &email.Sender{
			Email: email.Email{
				Address: o.InvitedByUser.OauthUser.Username,
				Name:    o.InvitedByUser.GetName(),
			},
		},
		Text: plainTextContent,
		HTML: htmlContent,
	}, nil
}

func newInvitationEmailPlainTextContent(title, name, company, invitedBy, link string) (string, error) {
	templateContents, err := ioutil.ReadFile(invitationEmailTemplateTxt)
	if err != nil {
		return "", err
	}
	tmpl, err := textTemplate.New(invitationEmailTemplateTxt).Parse(string(templateContents))
	if err != nil {
		return "", err
	}

	inventory := map[string]interface{}{
		"title":         title,
		"name":          name,
		"company":       company,
		"invitedBy":     invitedBy,
		"invitationURL": link,
	}

	var parsedTemplate bytes.Buffer
	if err := tmpl.Execute(&parsedTemplate, inventory); err != nil {
		return "", err
	}

	return parsedTemplate.String(), nil
}

func newInvitationEmailHTMLContent(title, inlineStyles, name, company, invitedBy, link string) (string, error) {
	// Layout
	layoutContents, err := ioutil.ReadFile(htmlEmailLayout)
	if err != nil {
		return "", err
	}
	layoutTmpl, err := htmlTemplate.New(htmlEmailLayout).Parse(string(layoutContents))
	if err != nil {
		return "", err
	}

	// Content
	templateContents, err := ioutil.ReadFile(invitationEmailTemplateHTML)
	if err != nil {
		return "", err
	}
	tmpl, err := htmlTemplate.New(invitationEmailTemplateHTML).Parse(string(templateContents))
	if err != nil {
		return "", err
	}

	var (
		inventory                   map[string]interface{}
		parsedContent, parsedLayout bytes.Buffer
	)

	// Parse the content template
	inventory = map[string]interface{}{
		"name":          name,
		"company":       company,
		"invitedBy":     invitedBy,
		"invitationURL": link,
	}
	if err := tmpl.Execute(&parsedContent, inventory); err != nil {
		return "", err
	}

	// Insert the content into the layout
	inventory = map[string]interface{}{
		"title":        title,
		"inlineStyles": htmlTemplate.CSS(inlineStyles),
		"content":      htmlTemplate.HTML(parsedContent.String()),
		"company":      company,
	}
	if err := layoutTmpl.Execute(&parsedLayout, inventory); err != nil {
		return "", err
	}

	return parsedLayout.String(), nil
}

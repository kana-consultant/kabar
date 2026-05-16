// services/email_service.go
package services

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"seo-backend/internal/config"
	"time"

	"github.com/go-mail/mail/v2"
)

type EmailService interface {
	SendInvitation(ctx context.Context, toEmail, token, teamID string) error
	SendWelcomeEmail(ctx context.Context, toEmail, name string) error
	SendPasswordReset(ctx context.Context, toEmail, token string) error
	SendEmail(ctx context.Context, to, subject, body string) error
}

type SMTPEmailService struct {
	config *config.SMTPConfig
	dialer *mail.Dialer
}

func NewSMTPEmailService(cfg *config.SMTPConfig) *SMTPEmailService {
	dialer := mail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password)

	// Configure TLS
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         cfg.Host,
	}

	// For Gmail and most providers
	if cfg.TLS {
		dialer.StartTLSPolicy = mail.MandatoryStartTLS
	}

	return &SMTPEmailService{
		config: cfg,
		dialer: dialer,
	}
}

// SendEmail - generic email sender
func (s *SMTPEmailService) SendEmail(ctx context.Context, to, subject, body string) error {
	m := mail.NewMessage()

	// Set sender
	from := s.config.FromEmail
	if s.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	}
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	// Set HTML body
	m.SetBody("text/html", body)

	// Add alternative plain text version (optional)
	// m.AddAlternative("text/plain", plainText)

	// Send email with timeout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Create a channel to handle send operation
	errChan := make(chan error, 1)
	go func() {
		errChan <- s.dialer.DialAndSend(m)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("email sending timeout: %w", ctx.Err())
	}
}

// SendInvitation - send team invitation email
func (s *SMTPEmailService) SendInvitation(ctx context.Context, toEmail, token, teamID string) error {
	inviteLink := fmt.Sprintf("%s/invite/?token=%s", getAppBaseURL(), token)

	subject := "Team Invitation"

	bodyTemplate := `
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
    </head>
    <body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
        <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
            <div style="background-color: #4F46E5; padding: 20px; text-align: center; border-radius: 8px 8px 0 0;">
                <h1 style="color: white; margin: 0;">Team Invitation</h1>
            </div>
            
            <div style="background-color: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; border: 1px solid #e5e7eb;">
                <p style="font-size: 16px;">Hello,</p>
                
                <p style="font-size: 16px;">You have been invited to join a team!</p>
                
                <div style="background-color: #ffffff; padding: 20px; border-radius: 8px; margin: 20px 0; border: 1px solid #e5e7eb;">
                    <p style="margin: 0 0 10px;"><strong>Team ID:</strong> {{.TeamID}}</p>
                    <p style="margin: 0;"><strong>Invitation Link:</strong></p>
                </div>
                
                <div style="text-align: center; margin: 30px 0;">
                    <a href="{{.InviteLink}}" 
                       style="background-color: #4F46E5; 
                              color: white; 
                              padding: 12px 30px; 
                              text-decoration: none; 
                              border-radius: 5px;
                              display: inline-block;">
                        Accept Invitation
                    </a>
                </div>
                
                <p style="font-size: 14px; color: #6b7280;">This invitation will expire in 7 days.</p>
                
                <hr style="margin: 20px 0; border: none; border-top: 1px solid #e5e7eb;">
                
                <p style="font-size: 12px; color: #9ca3af; text-align: center;">
                    If you didn't request this, you can safely ignore this email.
                </p>
            </div>
        </div>
    </body>
    </html>
    `

	tmpl, err := template.New("invitation").Parse(bodyTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	data := map[string]string{
		"InviteLink": inviteLink,
		"TeamID":     teamID,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	return s.SendEmail(ctx, toEmail, subject, body.String())
}

// SendWelcomeEmail - send welcome email after registration
func (s *SMTPEmailService) SendWelcomeEmail(ctx context.Context, toEmail, name string) error {
	subject := "Welcome to Our Platform!"

	bodyTemplate := `
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
    </head>
    <body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
        <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
            <div style="background-color: #10B981; padding: 20px; text-align: center; border-radius: 8px 8px 0 0;">
                <h1 style="color: white; margin: 0;">Welcome!</h1>
            </div>
            
            <div style="background-color: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; border: 1px solid #e5e7eb;">
                <p style="font-size: 16px;">Hi {{.Name}},</p>
                
                <p style="font-size: 16px;">Thank you for joining us! We're excited to have you on board.</p>
                
                <div style="text-align: center; margin: 30px 0;">
                    <a href="{{.DashboardURL}}" 
                       style="background-color: #10B981; 
                              color: white; 
                              padding: 12px 30px; 
                              text-decoration: none; 
                              border-radius: 5px;
                              display: inline-block;">
                        Go to Dashboard
                    </a>
                </div>
                
                <p style="font-size: 14px; color: #6b7280;">
                    Get started by exploring our features and setting up your account.
                </p>
            </div>
        </div>
    </body>
    </html>
    `

	tmpl, err := template.New("welcome").Parse(bodyTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	data := map[string]string{
		"Name":         name,
		"DashboardURL": getAppBaseURL() + "/dashboard",
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	return s.SendEmail(ctx, toEmail, subject, body.String())
}

// SendPasswordReset - send password reset email
func (s *SMTPEmailService) SendPasswordReset(ctx context.Context, toEmail, token string) error {
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", getAppBaseURL(), token)

	subject := "Password Reset Request"

	bodyTemplate := `
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
    </head>
    <body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
        <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
            <div style="background-color: #EF4444; padding: 20px; text-align: center; border-radius: 8px 8px 0 0;">
                <h1 style="color: white; margin: 0;">Reset Your Password</h1>
            </div>
            
            <div style="background-color: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; border: 1px solid #e5e7eb;">
                <p style="font-size: 16px;">Hello,</p>
                
                <p style="font-size: 16px;">We received a request to reset your password. Click the button below to create a new password:</p>
                
                <div style="text-align: center; margin: 30px 0;">
                    <a href="{{.ResetLink}}" 
                       style="background-color: #EF4444; 
                              color: white; 
                              padding: 12px 30px; 
                              text-decoration: none; 
                              border-radius: 5px;
                              display: inline-block;">
                        Reset Password
                    </a>
                </div>
                
                <p style="font-size: 14px; color: #6b7280;">This link will expire in 1 hour.</p>
                
                <hr style="margin: 20px 0; border: none; border-top: 1px solid #e5e7eb;">
                
                <p style="font-size: 12px; color: #9ca3af; text-align: center;">
                    If you didn't request this, please ignore this email.
                </p>
            </div>
        </div>
    </body>
    </html>
    `

	tmpl, err := template.New("reset-password").Parse(bodyTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	data := map[string]string{
		"ResetLink": resetLink,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	return s.SendEmail(ctx, toEmail, subject, body.String())
}

// Helper function to get app base URL
func getAppBaseURL() string {
	return getEnv("APP_BASE_URL", "http://localhost:3000")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

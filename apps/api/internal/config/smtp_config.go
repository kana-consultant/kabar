// config/smtp_config.go
package config

import (
	"strconv"
)

type SMTPConfig struct {
	Host      string
	Port      int
	User      string
	Password  string
	FromEmail string
	FromName  string
	Secure    bool
	TLS       bool
}

func LoadSMTPConfig() *SMTPConfig {
	port, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	secure, _ := strconv.ParseBool(getEnv("SMTP_SECURE", "false"))
	tls, _ := strconv.ParseBool(getEnv("SMTP_TLS", "true"))

	return &SMTPConfig{
		Host:      getEnv("SMTP_HOST", "smtp.gmail.com"),
		Port:      port,
		User:      getEnv("SMTP_USER", ""),
		Password:  getEnv("SMTP_PASSWORD", ""),
		FromEmail: getEnv("SMTP_FROM_EMAIL", ""),
		FromName:  getEnv("SMTP_FROM_NAME", "My App"),
		Secure:    secure,
		TLS:       tls,
	}
}

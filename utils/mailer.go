package utils

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

// EmailService - Struct chứa SMTP configuration
type EmailService struct {
	Host     string
	Port     int
	Username string
	Password string
	FromName string
}

// EmailData - Data để render email template
type EmailData struct {
	ApplicantName   string
	ApplicationCode string
	StatusText      string
	Note            string
}

// NewEmailService - Tạo instance EmailService từ env variables
func NewEmailService() *EmailService {
	return &EmailService{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     parseIntEnv("SMTP_PORT", 1025),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		FromName: os.Getenv("SMTP_FROM_NAME"),
	}
}

// SendApplicationStatusEmail - Gửi email thông báo status thay đổi
func (es *EmailService) SendApplicationStatusEmail(toEmail, applicantName, applicationCode, status, note string) error {
	// Validate email
	if toEmail == "" || !strings.Contains(toEmail, "@") {
		log.Printf("[Email Service] Invalid email address: %s", toEmail)
		return nil
	}

	statusText := getStatusDisplayText(status)

	data := EmailData{
		ApplicantName:   applicantName,
		ApplicationCode: applicationCode,
		StatusText:      statusText,
		Note:            note,
	}

	// Parse email template
	tmpl, err := template.ParseFiles("templates/email/status_update.html")
	if err != nil {
		log.Printf("[Email Service] Failed to parse template: %v", err)
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("[Email Service] Failed to execute template: %v", err)
		return err
	}

	// Build email message
	from := es.Username
	if from == "" {
		from = "noreply@localhost"
	}

	log.Printf("[Email Service] From address: %s", from)

	subject := "[Thông báo] Hồ sơ " + applicationCode + " - " + statusText
	body := buf.String()

	// Build MIME message
	message := "From: " + from + "\r\n" +
		"To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body

	// Connect to SMTP server
	addr := es.Host + ":" + strconv.Itoa(es.Port)
	log.Printf("[Email Service] Connecting to SMTP: %s, From: %s, To: %s", addr, from, toEmail)

	var auth smtp.Auth
	if es.Username != "" && es.Password != "" {
		auth = smtp.PlainAuth("", es.Username, es.Password, es.Host)
	}

	err = smtp.SendMail(addr, auth, from, []string{toEmail}, []byte(message))
	if err != nil {
		log.Printf("[Email Service] Failed to send email to %s: %v", toEmail, err)
		return err
	}

	log.Printf("[Email Service] ✓ Email sent successfully to %s for application %s", toEmail, applicationCode)
	return nil
}

// getStatusDisplayText - Convert status enum to display text
func getStatusDisplayText(status string) string {
	statusMap := map[string]string{
		"Received":            "Đã nhận",
		"Processing":          "Đang xử lý",
		"Supplement_Required": "Cần bổ sung",
		"Approved":            "Hoàn tất",
		"Rejected":            "Từ chối",
	}

	if text, ok := statusMap[status]; ok {
		return text
	}
	return status
}

// Helper function để parse int env
func parseIntEnv(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	result, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return result
}

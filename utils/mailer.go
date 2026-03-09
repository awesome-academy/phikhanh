package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/mail"
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
	FromAddr string
}

// EmailData - Data để render email template
type EmailData struct {
	ApplicantName   string
	ApplicationCode string
	StatusText      string
	Note            string
}

// WelcomeEmailData - Data để render welcome email template
type WelcomeEmailData struct {
	Name      string
	CitizenID string
	Email     string
	Password  string
	Role      string
	LoginURL  string
}

// NewEmailService - Tạo instance EmailService từ env variables
// Validate required config, log warning nếu config không đủ
func NewEmailService() *EmailService {
	es := &EmailService{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     parseIntEnv("SMTP_PORT", 1025),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		FromName: os.Getenv("SMTP_FROM_NAME"),
		FromAddr: getFromAddress(),
	}

	// Validate required SMTP configuration
	if err := es.ValidateConfig(); err != nil {
		log.Printf("[Email Service] Configuration Error: %v", err)
		log.Printf("[Email Service] Email notifications will be DISABLED")
	}

	return es
}

// ValidateConfig - Validate SMTP configuration có đủ parameters không
func (es *EmailService) ValidateConfig() error {
	if es.Host == "" {
		return fmt.Errorf("SMTP_HOST is required but empty")
	}

	if es.Port <= 0 || es.Port > 65535 {
		return fmt.Errorf("SMTP_PORT is invalid: %d (must be 1-65535)", es.Port)
	}

	if es.FromAddr == "" {
		return fmt.Errorf("SMTP_FROM_EMAIL (or SMTP_USERNAME) is required but empty")
	}

	// Log warning nếu authentication không được setup
	if es.Username == "" || es.Password == "" {
		log.Printf("[Email Service] Warning: SMTP authentication not configured (SMTP_USERNAME or SMTP_PASSWORD empty)")
	}

	return nil
}

// IsConfigured - Check nếu SMTP được configure đủ
func (es *EmailService) IsConfigured() bool {
	return es.Host != "" && es.Port > 0 && es.Port <= 65535 && es.FromAddr != ""
}

// getFromAddress - Get sender email address từ env hoặc default
func getFromAddress() string {
	addr := os.Getenv("SMTP_FROM_EMAIL")
	if addr != "" {
		return addr
	}

	// Fallback: dùng SMTP_USERNAME nếu là email
	username := os.Getenv("SMTP_USERNAME")
	if username != "" && strings.Contains(username, "@") {
		return username
	}

	return ""
}

// formatFromHeader - Format From header với display name
// Format: "FromName <from@example.com>" hoặc chỉ "from@example.com"
func (es *EmailService) formatFromHeader() string {
	if es.FromName != "" {
		// Proper RFC 2047 encoding nếu FromName chứa special characters
		address := mail.Address{
			Name:    es.FromName,
			Address: es.FromAddr,
		}
		return address.String()
	}
	return es.FromAddr
}

// SendApplicationStatusEmail - Gửi email thông báo status thay đổi
func (es *EmailService) SendApplicationStatusEmail(toEmail, applicantName, applicationCode, status, note string) error {
	// Kiểm tra SMTP được configure hay không
	if !es.IsConfigured() {
		log.Printf("[Email Service] Skipping email (SMTP not configured): to=%s, app=%s", toEmail, applicationCode)
		return NewBadRequestError("Email service not configured")
	}

	// Validate recipient email address
	if toEmail == "" {
		err := NewBadRequestError("Recipient email address is empty")
		log.Printf("[Email Service] %v for application %s", err, applicationCode)
		return err
	}

	if !strings.Contains(toEmail, "@") {
		err := NewBadRequestError(fmt.Sprintf("Invalid email address: %s", toEmail))
		log.Printf("[Email Service] %v for application %s", err, applicationCode)
		return err
	}

	// Validate sender email address
	if _, err := mail.ParseAddress(es.FromAddr); err != nil {
		err := NewBadRequestError(fmt.Sprintf("Invalid sender address: %s", es.FromAddr))
		log.Printf("[Email Service] %v", err)
		return err
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
		return NewInternalServerError(err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("[Email Service] Failed to execute template: %v", err)
		return NewInternalServerError(err)
	}

	// Format From header với display name
	fromHeader := es.formatFromHeader()
	subject := "[Thông báo] Hồ sơ " + applicationCode + " - " + statusText
	body := buf.String()

	// Build MIME message
	message := "From: " + fromHeader + "\r\n" +
		"To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body

	// Connect to SMTP server
	addr := es.Host + ":" + strconv.Itoa(es.Port)
	log.Printf("[Email Service] Connecting to SMTP: %s, From: %s, To: %s", addr, fromHeader, toEmail)

	var auth smtp.Auth
	if es.Username != "" && es.Password != "" {
		auth = smtp.PlainAuth("", es.Username, es.Password, es.Host)
	}

	err = smtp.SendMail(addr, auth, es.FromAddr, []string{toEmail}, []byte(message))
	if err != nil {
		log.Printf("[Email Service] Failed to send email to %s: %v", toEmail, err)
		return NewInternalServerError(err)
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

// parseIntEnv - Parse int từ env variable
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

// SendWelcomeEmail - Gửi welcome email cho user mới được tạo
func (es *EmailService) SendWelcomeEmail(user interface {
	GetEmail() string
	GetName() string
	GetCitizenID() string
	GetRole() string
}, rawPassword string) error {
	return es.sendWelcomeEmailTo(user.GetEmail(), user.GetName(), user.GetCitizenID(), user.GetRole(), rawPassword)
}

// sendWelcomeEmailTo - Internal helper để gửi welcome email
func (es *EmailService) sendWelcomeEmailTo(toEmail, name, citizenID, role, rawPassword string) error {
	if !es.IsConfigured() {
		log.Printf("[Email Service] Skipping welcome email (SMTP not configured): to=%s", toEmail)
		return NewBadRequestError("Email service not configured")
	}

	if toEmail == "" {
		return NewBadRequestError("Recipient email address is empty")
	}

	if !strings.Contains(toEmail, "@") {
		return NewBadRequestError(fmt.Sprintf("Invalid email address: %s", toEmail))
	}

	data := WelcomeEmailData{
		Name:      name,
		CitizenID: citizenID,
		Email:     toEmail,
		Password:  rawPassword,
		Role:      role,
	}

	// Parse email template từ file (giống status_update)
	tmpl, err := template.ParseFiles("templates/email/welcome.html")
	if err != nil {
		log.Printf("[Email Service] Failed to parse welcome template: %v", err)
		return NewInternalServerError(err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("[Email Service] Failed to execute welcome template: %v", err)
		return NewInternalServerError(err)
	}

	subject := "[Thông báo] Tài khoản của bạn đã được tạo"
	fromHeader := es.formatFromHeader()
	message := "From: " + fromHeader + "\r\n" +
		"To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		buf.String()

	addr := es.Host + ":" + strconv.Itoa(es.Port)

	var auth smtp.Auth
	if es.Username != "" && es.Password != "" {
		auth = smtp.PlainAuth("", es.Username, es.Password, es.Host)
	}

	if err := smtp.SendMail(addr, auth, es.FromAddr, []string{toEmail}, []byte(message)); err != nil {
		log.Printf("[Email Service] Failed to send welcome email to %s: %v", toEmail, err)
		return NewInternalServerError(err)
	}

	log.Printf("[Email Service] ✓ Welcome email sent to %s", toEmail)
	return nil
}

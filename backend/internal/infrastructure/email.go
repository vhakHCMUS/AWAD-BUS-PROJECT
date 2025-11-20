package infrastructure

import (
	"crypto/tls"
	"fmt"
	"os"

	gomail "gopkg.in/gomail.v2"
)

// EmailService handles email operations
type EmailService struct {
	dialer *gomail.Dialer
	from   string
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	host := os.Getenv("SMTP_HOST")
	port := 587
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	dialer := gomail.NewDialer(host, port, user, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false}

	return &EmailService{
		dialer: dialer,
		from:   from,
	}
}

// SendBookingConfirmation sends booking confirmation email
func (s *EmailService) SendBookingConfirmation(to, bookingCode string, attachmentPath string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Bus Booking Confirmation - "+bookingCode)

	body := fmt.Sprintf(`
		<h2>Booking Confirmed!</h2>
		<p>Your booking has been confirmed.</p>
		<p>Booking Code: <strong>%s</strong></p>
		<p>Please find your e-ticket attached.</p>
	`, bookingCode)

	m.SetBody("text/html", body)

	if attachmentPath != "" {
		m.Attach(attachmentPath)
	}

	return s.dialer.DialAndSend(m)
}

// SendTicket sends e-ticket via email
func (s *EmailService) SendTicket(to, ticketCode, pdfPath string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your Bus E-Ticket - "+ticketCode)

	body := `
		<h2>Your E-Ticket</h2>
		<p>Thank you for booking with us!</p>
		<p>Your e-ticket is attached to this email.</p>
		<p>Please present the QR code at the boarding point.</p>
	`

	m.SetBody("text/html", body)
	m.Attach(pdfPath)

	return s.dialer.DialAndSend(m)
}

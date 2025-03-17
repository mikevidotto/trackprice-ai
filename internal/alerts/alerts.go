package alerts

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// sendPriceAlert sends an email when a competitor's price changes
func SendPriceAlert(url, oldPrice, newPrice string) {
	// Get email config from .env
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	recipient := os.Getenv("ALERT_EMAIL")

	// Email body
	subject := "Price Change Alert: " + url
	body := fmt.Sprintf("A price change was detected for %s!\n\nOld Price: %s\nNew Price: %s", url, oldPrice, newPrice)
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	// Send email
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{recipient}, []byte(message))
	if err != nil {
		log.Println("❌ Failed to send price alert:", err)
	} else {
		fmt.Println("✅ Price change alert sent to", recipient)
	}
}

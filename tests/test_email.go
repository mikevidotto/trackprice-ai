package main

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func test_email() {
	// ✅ Load SendGrid API Key
	sendgridKey := os.Getenv("SENDGRID_API_KEY")
	if sendgridKey == "" {
		fmt.Println("❌ SendGrid API key is missing")
		return
	}

	// ✅ Define sender & recipient
	from := mail.NewEmail("TrackPrice AI", "no-reply@trackpriceai.com") // Change if needed
	to := mail.NewEmail("mike vidotto", "mikevidotto@live.com")         // Change to your real email

	// ✅ Email subject & content
	subject := "🚀 SendGrid Test Email from TrackPrice AI"
	content := "This is a test email from TrackPrice AI to verify SendGrid setup."

	// ✅ Create the email message
	message := mail.NewSingleEmail(from, subject, to, content, content)

	// ✅ Send the email
	client := sendgrid.NewSendClient(sendgridKey)
	response, err := client.Send(message)
	if err != nil {
		fmt.Printf("❌ Failed to send email: %v\n", err)
		return
	}

	// ✅ Print response status
	fmt.Printf("✅ Email sent! Status Code: %d\n", response.StatusCode)
}

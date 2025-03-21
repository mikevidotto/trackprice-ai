package notifications

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// ✅ SendPriceChangeAlert sends an email notification
func SendPriceChangeAlert(userEmail, competitorURL, oldPrice, newPrice string) error {
	sendgridKey := os.Getenv("SENDGRID_API_KEY")
	if sendgridKey == "" {
		return fmt.Errorf("❌ SendGrid API key is missing")
	}

	from := mail.NewEmail("TrackPrice AI", "no-reply@trackprice.ai") // ✅ Change to your domain if authenticated
	to := mail.NewEmail(userEmail, userEmail)
	subject := "🔔 Price Change Alert!"
	content := fmt.Sprintf(
		"Price change detected for %s!\n\nOld Price: %s\nNew Price: %s\n\nTrack more competitors with TrackPrice AI!",
		competitorURL, oldPrice, newPrice,
	)
	message := mail.NewSingleEmail(from, subject, to, content, content)

	client := sendgrid.NewSendClient(sendgridKey)
	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("❌ Failed to send email: %v", err)
	}

	fmt.Printf("✅ Email sent! Status Code: %d\n", response.StatusCode)
	return nil
}

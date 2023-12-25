package models

const (
	UserEventTypeEmailVerification = "user_verify_email"
)

const (
	EmailSubjectEmailVerification = "Verify email"
)

const (
	EmailBodyEmailVerification = `
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				.button {
					background-color: #007bff; /* Blue background */
					border: none;
					color: white;
					padding: 15px 32px;
					text-align: center;
					text-decoration: none;
					display: inline-block;
					font-size: 16px;
					margin: 4px 2px;
					cursor: pointer;
					border-radius: 5px;
				}
			</style>
		</head>
		<body>
		
			<p>Thank you for registering.</p>
			<p>One last step left - please verify your email address by clicking the button below.</p>
		
			<a href="%s" class="button">Verify</a>

		</body>
		</html>
	`
)

type UserMailItem struct {
	UserEventType string   `json:"user_event_type"`
	Receivers     []string `json:"receivers"`
	Link          string   `json:"link"`
}

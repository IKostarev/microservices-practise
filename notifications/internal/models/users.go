package models

const (
	UserEventTypeEmailVerification = "user_verify_email"
	UserEventTypeResetPassword     = "user_reset_password"
)

const (
	EmailSubjectEmailVerification = "Verify email"
	EmailSubjectResetPassword     = "Reset password"
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

	EmailBodyResetPassword = `
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

			<p>Someone requested password reset for your account. </p>
			<p>If it was not you - just ignore message.</p> 
			<p>If it was you - click the button</p>	

			<a href="%s" class="button">Reset</a>

		</body>
		</html>
	`
)

type UserMailItem struct {
	UserEventType string
	Receivers     []string
	Link          string
}

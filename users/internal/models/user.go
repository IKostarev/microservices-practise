package models

// UserDAO - data access object - струтктура для работы с базой данных
type UserDAO struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
}

// UserDTO - data transfer object - общая струтктура для передачи данных пользователя
type UserDTO struct {
	ID       int    `json:"id,omitempty" example:"1"`
	Username string `json:"username" example:"username"`
	Password string `json:"password,omitempty" example:"password"`
	Email    string `json:"email,omitempty" example:"user@example.com"`
}

// CreateUserDTO - data transfer object - струтктура для передачи данных пользователя при создании
type CreateUserDTO struct {
	ID                   int    `json:"id,omitempty" example:"1"`
	Username             string `json:"username" example:"username"`
	Password             string `json:"password,omitempty" example:"password"`
	PasswordConfirmation string `json:"password_confirmation,omitempty" example:"password"`
	Email                string `json:"email,omitempty" example:"user@example.com"`
}

// UpdateUserPasswordDTO - data transfer object - струтктура для передачи данных пользователя при обновлении пароля
type UpdateUserPasswordDTO struct {
	ID                   int    `json:"id,omitempty" example:"1"`
	OldPassword          string `json:"old_password" example:"password"`
	Password             string `json:"password" example:"password"`
	PasswordConfirmation string `json:"password_confirmation" example:"password"`
}

// UserLoginDTO - data transfer object - струтктура для передачи данных пользователя при логине
type UserLoginDTO struct {
	Username string `db:"username,omitempty"`
	Password string `db:"password"`
	Email    string `db:"email,omitempty"`
}

// UserTokens - струтктура для передачи токенов пользователя
type UserTokens struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

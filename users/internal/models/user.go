package models

import "users/pkg/grpc_stubs/users"

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

func NewEmptyUserDTO() *UserDTO {
	return &UserDTO{}
}

func (d *UserDTO) ToGRPC() *users.UserDTO {
	return &users.UserDTO{
		Id:       int32(d.ID),
		Username: d.Username,
		Password: d.Password,
		Email:    d.Email,
	}
}

func (d *UserDTO) FromGRPC(in *users.UserDTO) *UserDTO {
	d.ID = int(in.Id)
	d.Email = in.Email
	d.Username = in.Username
	d.Password = in.Password
	return d
}

// CreateUserDTO - data transfer object - струтктура для передачи данных пользователя при создании
type CreateUserDTO struct {
	ID                   int    `json:"id,omitempty" example:"1"`
	Username             string `json:"username" example:"username"`
	Password             string `json:"password,omitempty" example:"password"`
	PasswordConfirmation string `json:"password_confirmation,omitempty" example:"password"`
	Email                string `json:"email,omitempty" example:"user@example.com"`
}

func NewEmptyCreateUserDTO() *CreateUserDTO {
	return &CreateUserDTO{}
}

func (d *CreateUserDTO) ToGRPC() *users.CreateUserDTO {
	return &users.CreateUserDTO{
		Id:                   int32(d.ID),
		Username:             d.Username,
		Password:             d.Password,
		PasswordConfirmation: d.PasswordConfirmation,
		Email:                d.Email,
	}
}

func (d *CreateUserDTO) FromGRPC(in *users.CreateUserDTO) *CreateUserDTO {
	d.ID = int(in.Id)
	d.Email = in.Email
	d.Username = in.Username
	d.Password = in.Password
	d.PasswordConfirmation = in.PasswordConfirmation
	return d
}

// UpdateUserPasswordDTO - data transfer object - струтктура для передачи данных пользователя при обновлении пароля
type UpdateUserPasswordDTO struct {
	ID                   int    `json:"id,omitempty" example:"1"`
	OldPassword          string `json:"old_password" example:"password"`
	Password             string `json:"password" example:"password"`
	PasswordConfirmation string `json:"password_confirmation" example:"password"`
}

func NewEmptyUpdateUserPasswordDTO() *UpdateUserPasswordDTO {
	return &UpdateUserPasswordDTO{}
}

func (d *UpdateUserPasswordDTO) ToGRPC() *users.UpdateUserPasswordDTO {
	return &users.UpdateUserPasswordDTO{
		Id:                   int32(d.ID),
		OldPassword:          d.OldPassword,
		Password:             d.Password,
		PasswordConfirmation: d.PasswordConfirmation,
	}
}

func (d *UpdateUserPasswordDTO) FromGRPC(in *users.UpdateUserPasswordDTO) *UpdateUserPasswordDTO {
	d.ID = int(in.Id)
	d.OldPassword = in.OldPassword
	d.Password = in.Password
	d.PasswordConfirmation = in.PasswordConfirmation
	return d
}

// UserLoginDTO - data transfer object - струтктура для передачи данных пользователя при логине
type UserLoginDTO struct {
	Username string `db:"username,omitempty"`
	Password string `db:"password"`
	Email    string `db:"email,omitempty"`
}

func NewEmptyUserLoginDTO() *UserLoginDTO {
	return &UserLoginDTO{}
}

func (d *UserLoginDTO) ToGRPC() *users.UserLoginDTO {
	return &users.UserLoginDTO{
		Email:    d.Email,
		Username: d.Username,
		Password: d.Password,
	}
}

func (d *UserLoginDTO) FromGRPC(in *users.UserLoginDTO) *UserLoginDTO {
	d.Email = in.Email
	d.Username = in.Username
	d.Password = in.Password
	return d
}

// UserTokens - струтктура для передачи токенов пользователя
type UserTokens struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

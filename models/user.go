package models

import (
	"time"

	"goapp/adapter/gorm"
)

type Users []*User
type User struct {
	gorm.ModelUID
	Email              string
	Password           string
	IsAdmin            bool
	NewPasswordRequest *time.Time
}

type UserDtos []*UserDto
type UserDto struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

type UserLoginForm struct {
	Email     string `json:"email" form:"required,max=255,email"`
	Password  string `json:"password"  form:"required"`
	UseCookie bool   `json:"use_cookie,omitempty"`
}

type UserRegisterEmailFormDto struct {
	Email string `json:"email"`
}
type UserRegisterEmailForm struct {
	Email string    `json:"email" form:"required,max=255,email"`
	Date  time.Time `json:"date"`
}

type UserRegisterPasswordForm struct {
	Password string `json:"password" form:"required,min=6,max=100"`
}

func (f *UserLoginForm) ToModel() (*User, error) {
	return &User{
		Email:    f.Email,
		Password: f.Password,
	}, nil
}

func (f *UserRegisterEmailForm) ToModel() (*User, error) {
	return &User{
		Email: f.Email,
	}, nil
}

func (u UserRegisterEmailForm) ToDto() *UserRegisterEmailFormDto {
	return &UserRegisterEmailFormDto{
		Email: u.Email,
	}
}
func (u User) ToDto() *UserDto {
	return &UserDto{
		ID:      u.ID,
		Email:   u.Email,
		IsAdmin: u.IsAdmin,
	}
}

/*
// An User Object
type User struct {
	Auth            int     `json:"auth" example:"1"`
	Oid             string  `json:"oid" example:"cidb"`
	Name            string  `json:"name" example:"Mustermann Max"`
	Email           string  `json:"email" example:"max.mustermann@demo-mandant.de"`
	Bezeichnung     string  `json:"bezeichnung" example:"Max Mustermann"`
	Token           string  `json:"token" example:"8C08E916 ... D4349C59"`
	JWTToken        *string `json:"jwtToken,omitempty" example:"eyJ0e ... gNwo5bS0v9qc"`
	JwtRefreshToken *string `json:"jwtRefreshToken,omitempty" example:"eyJ0e ... gNwo5bS0v9qc"`
}
*/

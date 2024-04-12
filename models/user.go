package models

import (
	"time"

	g "goapp/adapter/gorm"
	"goapp/util/tools"
)

type Users []*User
type User struct {
	g.ModelUID
	Email              string
	Password           string
	Newsletter         bool
	FirstName          string
	LastName           string
	IsAdmin            bool
	NewPasswordRequest *time.Time
}
type RegisterUser struct {
	g.ModelUID
	Email      string
	Password   string
	Newsletter bool
	FirstName  string
	LastName   string
}
type UserRegisterResponse struct {
	Message string `json:"message" `
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
type UserPasswordForm struct {
	Link     string `json:"link" `
	Password string `json:"password" form:"required,min=5,max=64"`
}
type RegisterUserForm struct {
	Email      string `json:"email" form:"required,max=255,email"`
	FirstName  string `json:"first_name" form:"required,min=2,max=255"`
	LastName   string `json:"last_name" form:"required,min=2,max=255"`
	Newsletter bool   `json:"newsletter"`
	TermsAgree bool   `json:"terms_agree" form:"required"`
	Password   string `json:"password" form:"required,min=5,max=64"`
}

type RegisterUserLink struct {
	ID        string
	CreatedAt time.Time
}

type UserPasswordLink struct {
	ID                 string
	NewPasswordRequest time.Time
}

func (u User) ToDto() *UserDto {
	return &UserDto{
		ID:      u.ID,
		Email:   u.Email,
		IsAdmin: u.IsAdmin,
	}
}

func (f *RegisterUserForm) ToModel() (*RegisterUser, error) {

	hashedPassword, err := tools.HashPassword(f.Password)

	if err != nil {
		return nil, err
	}

	return &RegisterUser{
		Email:      f.Email,
		Password:   hashedPassword,
		FirstName:  f.FirstName,
		LastName:   f.LastName,
		Newsletter: f.Newsletter,
	}, nil
}

func (f *RegisterUser) ToLinkModel() *RegisterUserLink {

	user := &RegisterUserLink{
		ID:        f.ID,
		CreatedAt: f.CreatedAt,
	}
	return user
}
func (f *RegisterUser) ToUserModel() *User {

	user := &User{
		Email:     f.Email,
		Password:  f.Password,
		FirstName: f.FirstName,
		LastName:  f.LastName,
	}
	return user
}

func (f *RegisterUserLink) ToModel() *RegisterUser {

	user := &RegisterUser{}
	user.ID = f.ID
	user.CreatedAt = f.CreatedAt

	return user
}

/*

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



*/
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

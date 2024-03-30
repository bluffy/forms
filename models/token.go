package models

type Token struct {
	AccessToken  string `json:"at"`
	RefreshToken string `json:"rt"`
}

/*
type TokenDto struct {
	//UID          string `json:"id"`
	Email        string `json:"email"`
	IsAdmin      bool   `json:"isAdmin"`
	AccessToken  string `json:"at"`
	RefreshToken string `json:"rt"`
}

func (t Token) ToDto(u *User) *TokenDto {
	return &TokenDto{
		UID:          u.ModelUID.ID,
		Email:        u.Email,
		IsAdmin:      u.IsAdmin,
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	}
}
*/

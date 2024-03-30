package models

type Token struct {
	AccessToken  string `json:"jwtToken"`
	RefreshToken string `json:"jwtRefreshToken"`
}

type TokenDto struct {
	Email        string `json:"email"`
	AccessToken  string `json:"jwtToken"`
	RefreshToken string `json:"jwtRefreshToken"`
}

package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bluffy/forms/models"
	"github.com/golang-jwt/jwt"
)

type Jwt struct {
	TokenLifeTime        int
	TokenRefreshLifeTime int
	TokenRefreshAllowd   bool
	TokenKey             string
}

type JWTClaim struct {
	Data models.UserDto
	jwt.StandardClaims
}

func (j Jwt) CreateToken(user models.UserDto) (models.Token, error) {
	var err error

	claims := &JWTClaim{
		Data: user,
	}
	//claims.ExpiresAt = time.Now().Add(time.Hour * 12).Unix()
	//claims.ExpiresAt = time.Now().Add(time.Second * 5).Unix()
	claims.ExpiresAt = time.Now().Add(time.Minute * time.Duration(j.TokenLifeTime)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt := models.Token{}

	jwt.AccessToken, err = token.SignedString([]byte(j.TokenKey))
	if err != nil {
		return jwt, err
	}

	return j.createRefreshToken(jwt)
}

func (j Jwt) ValidateToken(accessToken string) (models.UserDto, error) {

	token, err := jwt.ParseWithClaims(
		accessToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(j.TokenKey), nil
		},
	)
	user := models.UserDto{}

	if err != nil {
		return user, err
	}

	payload, ok := token.Claims.(*JWTClaim)
	if ok && token.Valid {
		user = payload.Data
		return user, nil
	}

	return user, errors.New("invalid token")
}

func (j Jwt) ValidateRefreshToken(modelToken models.Token) (models.UserDto, error) {

	user := models.UserDto{}

	refreshToken, err := jwt.Parse(modelToken.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.TokenKey), nil
	})
	if err != nil {
		return user, err
	}

	refreshPayload, refeshOk := refreshToken.Claims.(jwt.MapClaims)
	if !(refeshOk && refreshToken.Valid) {
		return user, errors.New("invalid refreshToken")
	}

	token, err := jwt.ParseWithClaims(
		refreshPayload["token"].(string),
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(j.TokenKey), nil
		},
	)

	// expirest token is allowed
	if err != nil {
		if !strings.HasPrefix(err.Error(), "token is expired by ") {
			return user, err
		}
	}

	payload, ok := token.Claims.(*JWTClaim)

	if ok {
		user = payload.Data

		//FEhlt noch, check token with DATABASE

		return user, nil
	}

	return user, errors.New("invalid token")

}

func (j Jwt) createRefreshToken(token models.Token) (models.Token, error) {
	var err error

	claims := jwt.MapClaims{}
	claims["token"] = token.AccessToken
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(j.TokenRefreshLifeTime)).Unix()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token.RefreshToken, err = refreshToken.SignedString([]byte(j.TokenKey))
	if err != nil {
		return token, err
	}

	return token, nil
}

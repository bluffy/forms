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

type JWTClaimData struct {
	User      models.UserDto
	SessionId string
}

type JWTClaim struct {
	Data JWTClaimData
	jwt.StandardClaims
}

func (j Jwt) CreateToken(user models.UserDto, sessionID string) (models.Token, error) {
	var err error

	claimData := &JWTClaimData{
		User:      user,
		SessionId: sessionID,
	}
	claims := &JWTClaim{
		Data: *claimData,
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

func (j Jwt) ValidateToken(accessToken string) (*models.UserDto, *string, error) {

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
		return nil, nil, err
	}

	payload, ok := token.Claims.(*JWTClaim)
	if ok && token.Valid {
		user = payload.Data.User
		sessionId := payload.Data.SessionId
		return &user, &sessionId, nil
	}

	return nil, nil, errors.New("invalid token")
}

func (j Jwt) ValidateRefreshToken(modelToken models.Token) (*models.UserDto, *string, error) {

	refreshToken, err := jwt.Parse(modelToken.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.TokenKey), nil
	})
	if err != nil {
		return nil, nil, err
	}

	refreshPayload, refeshOk := refreshToken.Claims.(jwt.MapClaims)
	if !(refeshOk && refreshToken.Valid) {
		return nil, nil, errors.New("invalid refreshToken")
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
			return nil, nil, err
		}
	}

	payload, ok := token.Claims.(*JWTClaim)

	if ok {
		user := payload.Data.User
		sessionId := payload.Data.SessionId

		//FEhlt noch, check token with DATABASE

		return &user, &sessionId, nil
	}

	return nil, nil, errors.New("invalid token")

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

package app

import (
	"encoding/json"
	"net/http"
	"time"

	"goapp/config"
	"goapp/models"
	"goapp/repository"
	"goapp/service"
	"goapp/util/tools"

	log "github.com/sirupsen/logrus"
)

// HandlerLogin  godoc
// @Tags         public
// @Description  login
// @Accept       json
// @Produce      json
// @Param data body models.UserLoginForm  true "Email & Password"
// @Success      200 {object} models.Token
// @Failure      422 {object} models.AppError
// @Failure      500 {object} models.AppError
// @Router       /api/v1/login [post]
func (app *App) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	form := &models.UserLoginForm{}
	if app.checkForm(form, w, r) {
		return
	}

	user, err := repository.GetUserByEmail(app.db, form.Email)
	if err != nil {
		app.printError(w, http.StatusUnprocessableEntity, 200, err, "")
		return
	}
	if !tools.CheckPasswordHash(form.Password, user.Password) {
		app.printError(w, http.StatusUnprocessableEntity, 201, err, "")
		return
	}

	session := &models.Session{
		UserID: user.ID,
	}

	session, err = repository.CreateSession(app.db, session)
	if err != nil {
		app.printError(w, http.StatusInternalServerError, 104, err, "")
		return
	}
	sessionEnc, err := tools.EncryptBase64(session.ID, config.Conf.EncryptKey)
	if err != nil {
		app.printError(w, http.StatusInternalServerError, 104, err, "")
		return
	}

	jwt := service.Jwt{
		TokenLifeTime:        config.Conf.Server.TokenLifeTime,
		TokenRefreshLifeTime: config.Conf.Server.TokenRefreshLifeTime,
		TokenRefreshAllowd:   config.Conf.Server.TokenRefreshAllowed,
		TokenKey:             config.Conf.Server.TokenKey,
	}

	userDto := user.ToDto()
	token, err := jwt.CreateToken(*userDto, session.ID)
	if err != nil {
		app.printError(w, http.StatusInternalServerError, 202, err, "")
		return
	}
	/*

		at := http.Cookie{
			Name: "at",
			Path: "/",
			//Domain: "localhost",
			//MaxAge:   3600,
			HttpOnly: true,
			Expires:  time.Now().AddDate(1, 0, 0),
			SameSite: http.SameSiteLaxMode,
			Value: token.AccessToken,
		}
		rt := http.Cookie{
			Name: "at",
			Path: "/",
			//Domain: "localhost",
			//MaxAge:   3600,
			HttpOnly: true,
			Expires:  time.Now().AddDate(1, 0, 0),
			SameSite: http.SameSiteLaxMode,
			Value: token.AccessToken,
		}
	*/

	/*
		var sess session.Store

		sess.Set("session", "session middleware")
	*/

	//var Session session.Store

	sessionCookie := http.Cookie{
		Name:     "session",
		Path:     "/",
		HttpOnly: false,
		Expires:  time.Now().Add(time.Second * time.Duration(config.Conf.Server.TokenLifeTime)),
		//SameSite: http.SameSiteLaxMode,
		SameSite: http.SameSiteLaxMode,
		Value:    sessionEnc,
		//Secure:   true,
	}
	http.SetCookie(w, &sessionCookie)

	at := http.Cookie{
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Second * time.Duration(config.Conf.Server.TokenLifeTime)),
		SameSite: http.SameSiteLaxMode,
		Value:    token.AccessToken,
	}
	http.SetCookie(w, &at)

	if config.Conf.Server.TokenRefreshAllowed {

		rt := http.Cookie{
			Name:     "rt",
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(time.Second * time.Duration(config.Conf.Server.TokenRefreshLifeTime)),
			SameSite: http.SameSiteLaxMode,
			Value:    token.RefreshToken,
		}
		http.SetCookie(w, &rt)

	}
	//dtos := token.ToDto(user)

	if err := json.NewEncoder(w).Encode(token); err != nil {
		log.Warn(err)
		app.printError(w, http.StatusInternalServerError, 102, err, "")
	}

}

func (app *App) RefreshLoginToken(w http.ResponseWriter, r *http.Request) {

	token := models.Token{}

	log.Debug("body: %d", r.Context())
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		app.printError(w, http.StatusUnprocessableEntity, 103, err, "")
		return
	}

	jwt := service.Jwt{}
	user, sessionId, err := jwt.ValidateRefreshToken(token)
	if err != nil {
		app.printError(w, http.StatusUnprocessableEntity, 203, err, "")
		return
	}

	token, err = jwt.CreateToken(*user, *sessionId)
	if err != nil {
		app.printError(w, http.StatusInternalServerError, 204, err, "")
		return
	}

	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		app.printError(w, http.StatusInternalServerError, 102, err, "")
		return
	}
}

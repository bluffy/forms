package app

import (
	"embed"

	"github.com/bluffy/forms/config"
	"github.com/bluffy/forms/models"
	"github.com/bluffy/forms/server/lang"
	"gorm.io/gorm"
)

//go:embed version/*
var versionFS embed.FS
var version = "0.0.0"

type App struct {
	conf *config.Config
	lang *lang.Lang
	db   *gorm.DB
	//userRestConf *clientcredentials.Config
	//	openIds map[string]*oauth2.Config
	//userClient *http.Client
	//amsClient  *http.Client
	user *models.DatabaseUser
}

func New(
	conf *config.Config,
	lang *lang.Lang,
	db *gorm.DB,
) *App {

	return &App{
		conf: conf,
		lang: lang,
		//openIds: openIds,
		db: db,
	}
}

func init() {
	data, _ := versionFS.ReadFile("version/VERSION")
	if data != nil {
		version = string(data)
	}
}

func GetVersion() string {
	return version
}

func (app *App) SetUser(user models.DatabaseUser) {
	app.user = &user
}

func (app *App) User() *models.DatabaseUser {
	return app.user
}

func (app *App) Conf() *config.Config {
	return app.conf
}

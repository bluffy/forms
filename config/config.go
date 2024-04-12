package config

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Dev            bool   `default:"false"`
	Debug          bool   `default:"false"`
	ShowApiDoku    bool   `yaml:"show_api_doku" default:"false"`
	Language       string `default:"en"`
	UseEmbedClient bool   `yaml:"use_embed_client" default:"true"`
	EncryptKey     string `yaml:"encryptKey" default:"6GbE8Qf2rkbYm9EecnxfVnBzXp8ZvWo6h3FDKxA88qv46U8ueRY4RJcbD7oMjCAzQLT"`
	Smtp           Smtp   `yaml:"smtp"`
	Database       Database
	Server         ServerConf
}

type ServerConf struct {
	Port         int           `yaml:"port" default:"4090"`
	PortIntern   int           `yaml:"port_intern" default:"4091"`
	Host         string        `default:"localhost"`
	PublicURL    string        `yaml:"public_url" default:"http://localhost:4090"`
	ClientUrl    string        `yaml:"client_url"`
	TimeoutRead  time.Duration `yaml:"timeout_read" default:"default=30s"`
	TimeoutWrite time.Duration `yaml:"timeout_write" default:"default=30s"`
	TimeoutIdle  time.Duration `yaml:"timeout_idle" default:"default=45s"`
	Cors         Cors
}

type MysqlConf struct {
	Username string `default:"goapp"`
	Password string `default:"goapp"`
	Database string `default:"goapp"`
	Host     string `default:"localhost"`
	Port     string `default:"3306"`
}
type SqliteConf struct {
	Path string `default:"database.db"`
}

type Database struct {
	Type   string `default:"sqlite"`
	Mysql  MysqlConf
	Sqlite SqliteConf
}
type Smtp struct {
	Host             string
	User             string
	Password         string
	Port             int `default:"587"`
	Sender           string
	SkipSSLVerify    bool `yaml:"skip_ssl_verify" default:"true"`
	SendWithStartTLS bool `yaml:"send_with_start_tls" default:"true"`
	SendWithTLS      bool `yaml:"send_with_ssl" default:"false"`
}
type Cors struct {
	AllowedOrigins   []string `yaml:"allowed_orgins" default:"[\"https://*\",\"http://*\"]"`
	AllowCredentials bool     `yaml:"allowed_credential" default:"true"`
	AllowedMethods   []string `yaml:"allowed_methods" default:"[\"GET\",\"POST\",\"PUT\",\"DELETE\",\"OPTIONS\",\"HEAD\"]"`
	AllowedHeaders   []string `yaml:"allowed_headers" default:"[\"Accept\",\"Authorization\",\"Content-Type\",\"X-CSRF-Token\",\"X-Requested-With\"]"`
	ExposedHeaders   []string `yaml:"exposed_headers" default:"[\"Link\",\"set-cookie\"]}"`
	MaxAge           int      `yaml:"max_age" default:"300"`
}

func (conf *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(conf)
	//logrus.Infof("#######test:%+v\n", conf)
	type plain Config
	if err := unmarshal((*plain)(conf)); err != nil {
		return err
	}
	return nil
}

func New(configFile string) (*Config, error) {

	config := &Config{}

	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	logrus.Info("Config File: " + configFile)

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	if config.Server.PublicURL == "" {
		return nil, errors.New("config missing: server.public_url (https://domain.com/)")
	}
	config.Server.PublicURL = strings.TrimRight(config.Server.PublicURL, "/")

	if config.Server.ClientUrl == "" {
		config.Server.ClientUrl = config.Server.PublicURL
	}
	config.Server.ClientUrl = strings.TrimRight(config.Server.ClientUrl, "/")

	return config, nil
}

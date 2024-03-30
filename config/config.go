package config

import (
	"os"
	"time"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Dev         bool
	Debug       bool
	Language    string `default:"en"`
	LogLanguage string `yaml:"log_language" default:"en"`

	Server struct {
		Port                 int           `yaml:"port" default:"4090"`
		PortIntern           int           `yaml:"port_intern" default:"4091"`
		Host                 string        `default:"127.0.0.1"`
		PublicUrl            string        `yaml:"public_url"`
		TimeoutRead          time.Duration `yaml:"timeout_read" default:"default=30s"`
		TimeoutWrite         time.Duration `yaml:"timeout_write" default:"default=30s"`
		TimeoutIdle          time.Duration `yaml:"timeout_idle" default:"default=45s"`
		TokenKey             string        `yaml:"token_key"`
		TokenLifeTime        int           `yaml:"token_life_time" default:"43200"`
		TokenRefreshLifeTime int           `yaml:"token_refreshLifeTime" default:"43200"`
		TokenRefreshAllowed  bool          `yaml:"token_refreshAllowed" default:"true"`
	}

	Database struct {
		Type     string
		Username string
		Password string
		Database string
		Host     string
		Port     string
	}
}

func (conf *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(conf)

	type plain Config
	if err := unmarshal((*plain)(conf)); err != nil {
		return err
	}
	return nil
}

func AppConfig(configFile string) (*Config, error) {

	config := &Config{}

	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

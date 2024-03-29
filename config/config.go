package config

import (
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Dev      bool
	Debug    bool
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

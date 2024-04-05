package lang

import (
	"embed"
	"io/fs"

	"github.com/creasty/defaults"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"goapp/lang/locale"
	"path/filepath"
)

type Region struct {
	Text  map[string]string
	Error map[int]string
}

type Locale struct {
	Text      locale.Text      `yaml:"text"`
	Validator locale.Validator `yaml:"validator"`
	Error     locale.Error     `yaml:"error"`
}

type Lang struct {
	Locale        map[string]Locale
	DefaultLocale *Locale
	Default       string
	Log           string
}

func AppLang(defaultLang string, langFS embed.FS) *Lang {
	var files []string
	var language Lang
	language.Default = defaultLang
	language.Locale = make(map[string]Locale)

	locale := &Locale{}

	err := yaml.Unmarshal([]byte(``), &locale)
	if err != nil {
		log.Info(err)
		log.Fatal("Error in Example")
	}
	language.DefaultLocale = locale

	language.Locale[language.Default] = *locale

	if err := fs.WalkDir(langFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yaml" {
			return nil
		}

		files = append(files, path)

		return nil

	}); err != nil {
		log.Fatal("Error in Lang Walkdir")
		return nil
	}

	for _, file := range files {
		f2, _ := langFS.Open(file)
		//b, _ := io.ReadAll(f)

		d2 := yaml.NewDecoder(f2)

		locale := &Locale{}

		err = d2.Decode(&locale)
		if err != nil {

			log.Info("Error Lang Parsing: " + file)
			log.Info(err)
		} else {
			code := file[:len(file)-len(filepath.Ext(file))]
			code = filepath.Base(code)
			language.Locale[code] = *locale

			if code == language.Default {
				language.DefaultLocale = locale
			}

		}

	}

	_, ok := language.Locale[language.Default]
	if !ok {
		log.Fatal("Error in Lang: Default lang: " + language.Default + " not exists!")
		return nil
	}

	return &language
}

func (conf *Locale) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(conf)
	type plain Locale
	if err := unmarshal((*plain)(conf)); err != nil {
		return err
	}
	return nil
}

/*

func unmarshalYAML(file string, unmarshal func(interface{}) error) error {
	defaults.Set(file)

	type plain Config
	if err := unmarshal((*plain)(conf)); err != nil {
		return err
	}
	return nil
}
*/

/*
func (conf *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(conf)

	type plain Config
	if err := unmarshal((*plain)(conf)); err != nil {
		return err
	}
	return nil
}
*/
/*
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
*/

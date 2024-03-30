package lang

import (
	"embed"
	"io/fs"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"path/filepath"
)

type Region struct {
	Text  map[string]string
	Error map[int]string
}

type Lang struct {
	Region  map[string]Region
	Default string
	Log     string
}

/*
type Lang map[string]Code
type Code struct {
	Error map[string]string
}
*/

/*
	func init() {
		if _, err := toml.DecodeFile("./data/lang/de.toml", &Lang.Error); err != nil {
			panic(err)
		}
	}
	func AppLang(confLang string, langFS embed.FS) *Lang {
	language := make(Lang, 1)
	code := &Code{}
	if _, err := toml.DecodeFS(langFS, "data/lang/en.toml", &code); err != nil {
		panic(err)
	}

		if _, err := toml.DecodeFile("./data/lang/de.toml", &code); err != nil {
			panic(err)
		}

	language[confLang] = *code
	return &language
	}
*/

func AppLang(defaultLang string, logLang string, langFS embed.FS) *Lang {
	var files []string
	var language Lang
	language.Default = defaultLang
	language.Log = logLang
	language.Region = make(map[string]Region)

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
		f, _ := langFS.Open(file)
		//b, _ := io.ReadAll(f)

		d := yaml.NewDecoder(f)
		region := &Region{}

		err := d.Decode(&region)
		if err != nil {
			log.Info("Error Lang Parsing: " + file)
			log.Info(err)
		} else {
			code := file[:len(file)-len(filepath.Ext(file))]
			code = filepath.Base(code)
			language.Region[code] = *region
		}

	}

	_, ok := language.Region[language.Default]
	if !ok {
		log.Fatal("Error in Lang: Default lang: " + language.Default + " not exists!")
		return nil
	}

	_, ok = language.Region[language.Log]
	if !ok {
		log.Fatal("Error in Lang: Log lang: " + language.Log + " not exists!")
		return nil
	}

	return &language
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

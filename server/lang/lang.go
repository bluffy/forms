package lang

import (
	"embed"

	"github.com/BurntSushi/toml"
)

type Lang map[string]Code
type Code struct {
	Error map[string]string
}

/*
	func init() {
		if _, err := toml.DecodeFile("./data/lang/de.toml", &Lang.Error); err != nil {
			panic(err)
		}
	}
*/
func AppLang(confLang string, langFS embed.FS) *Lang {
	language := make(Lang, 1)
	code := &Code{}
	if _, err := toml.DecodeFS(langFS, "data/lang/en.toml", &code); err != nil {
		panic(err)
	}
	/*
		if _, err := toml.DecodeFile("./data/lang/de.toml", &code); err != nil {
			panic(err)
		}
	*/
	language[confLang] = *code
	return &language
}

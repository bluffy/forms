package locale

type Validator struct {
	Required    string `default:"is a required"`
	Max         string `default:"must be a maximum of %s in length"`
	Min         string `default:"must be a minimum of %s in length"`
	Url         string `default:"must be a valid URL"`
	Email       string `default:"must be a valid email address"`
	Alpha_space string `default:"can only contain alphabetic and space characters"`
	Date        string `default:"must be a valid date"`
	Default     string `default:"something wrong: %s"`
}

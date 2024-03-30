package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	alphaSpaceRegexString string = "^[a-zA-Z ]+$"
	dateRegexString       string = "^(((19|20)([2468][048]|[13579][26]|0[48])|2000)[/-]02[/-]29|((19|20)[0-9]{2}[/-](0[469]|11)[/-](0[1-9]|[12][0-9]|30)|(19|20)[0-9]{2}[/-](0[13578]|1[02])[/-](0[1-9]|[12][0-9]|3[01])|(19|20)[0-9]{2}[/-]02[/-](0[1-9]|1[0-9]|2[0-8])))$"
)

type ErrResponse struct {

	//Errors []string `json:"errors"`
	Error struct {
		Message string                 `json:"message,omitempty"`
		Fields  map[string]interface{} `json:"fields,omitempty"`
	} `json:"error"`
}

func isAlphaSpace(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(alphaSpaceRegexString)
	return reg.MatchString(fl.Field().String())
}
func isDate(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(dateRegexString)
	return reg.MatchString(fl.Field().String())
}

func New() *validator.Validate {
	validate := validator.New()
	validate.SetTagName("form")

	// Using the names which have been specified for JSON representations of structs, rather than normal Go field names

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	validate.RegisterValidation("alpha_space", isAlphaSpace)
	validate.RegisterValidation("date", isDate)

	return validate
}
func ErrorMsg(msg string) *ErrResponse {
	resp := ErrResponse{}
	resp.Error.Message = msg
	return &resp

}
func ToErrResponse(err error, msg *string) *ErrResponse {

	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		/*
			resp := ErrResponse{
				Errors: make(map[string]interface{}),
			}
		*/
		resp := ErrResponse{}
		resp.Error.Fields = make(map[string]interface{})
		if msg != nil {
			resp.Error.Message = *msg
		}

		for _, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				resp.Error.Fields[err.Field()] = fmt.Sprintf("%s is a required field", err.Field())
			case "max":
				resp.Error.Fields[err.Field()] = fmt.Sprintf("%s must be a maximum of %s in length", err.Field(), err.Param())
			case "min":
				resp.Error.Fields[err.Field()] = fmt.Sprintf("%s must be a minimum of %s in length", err.Field(), err.Param())
			case "url":
				resp.Error.Fields[err.Field()] = fmt.Sprintf("%s must be a valid URL", err.Field())
			case "email":
				resp.Error.Fields[err.Field()] = fmt.Sprintf("%s must be a valid email address", err.Field())
			case "alpha_space":
				resp.Error.Fields[err.Field()] = fmt.Sprintf("%s can only contain alphabetic and space characters", err.Field())
			case "date":
				resp.Error.Fields[err.Field()] = fmt.Sprintf("%s must be a valid date", err.Field())
			default:
				resp.Error.Fields[err.Field()] = fmt.Sprintf("something wrong on %s; %s", err.Field(), err.Tag())
			}
		}
		return &resp
	}
	return nil
}

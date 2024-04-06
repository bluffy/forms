package locale

type Text struct {
	Welcome string `default:"Welcome"`
	Page    struct {
		Auth struct {
			Register_link_is_invalid         string `default:"the link is invalid or expired"`
			Register_link_is_expired         string `default:"the link is expired"`
			Register_user_already_registered string `default:"user alread exists"`
			Message_user_created             string `default:"successful, you can login now!"`
		}
	}

	Error__commen_server_error string `default:"Error from Server"`
	Error__database_error      string `default:"Server Database Error"`
	Error__json_create         string `default:"Server Error: JSON Create"`
	Error__json_encode         string `default:"Server Error: Encode"`
	Error__json_decode         string `default:"Server Error: Decode"`
	Error__form_response_error string `default:"fomrs Error"`

	Session__error_sessen_not_exists_or_expired string `default:"session not exists or expired!"`

	Page_auth__regsitering_success            string `default:"register was successful, check your email please!"`
	Page_auth__error__user_not_exists         string `default:"user & password not matched"`
	Page_auth__error__wrong_password          string `default:"user & password not matched"`
	Page_auth__error__session_regeneration    string `default:"unable to create session"`
	Page_auth__error__user_already_registered string `default:"user already registered"`

	//Page_auth__error__invalid_token              string `default:"invalid session"`
	//Page_auth__error__unable_create_access_token string `default:"unable to create access token"`
}

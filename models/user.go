package models

// An User Object
type DatabaseUser struct {
	Auth            int     `json:"auth" example:"1"`
	Oid             string  `json:"oid" example:"cidb"`
	Name            string  `json:"name" example:"Mustermann Max"`
	Email           string  `json:"email" example:"max.mustermann@demo-mandant.de"`
	Bezeichnung     string  `json:"bezeichnung" example:"Max Mustermann"`
	Token           string  `json:"token" example:"8C08E916 ... D4349C59"`
	JWTToken        *string `json:"jwtToken,omitempty" example:"eyJ0e ... gNwo5bS0v9qc"`
	JwtRefreshToken *string `json:"jwtRefreshToken,omitempty" example:"eyJ0e ... gNwo5bS0v9qc"`
}

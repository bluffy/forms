// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "dev",
            "email": "dev@vocy.de"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/bl-api/page/v1/user/login": {
            "post": {
                "description": "login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "parameters": [
                    {
                        "description": "Email \u0026 Password",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserLoginForm"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/app.ErrResponse"
                        }
                    }
                }
            }
        },
        "/bl-api/page/v1/user/register": {
            "post": {
                "description": "login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "parameters": [
                    {
                        "description": "register informations",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.RegisterUserForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.ApiPageResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/app.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Response JSON",
                        "schema": {
                            "$ref": "#/definitions/app.ErrResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "app.ApiPageResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "app.ErrResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "description": "Errors []string ` + "`" + `json:\"errors\"` + "`" + `",
                    "type": "object",
                    "properties": {
                        "fields": {
                            "type": "object",
                            "additionalProperties": true
                        },
                        "message": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "models.RegisterUserForm": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "newsletter": {
                    "type": "boolean"
                },
                "password": {
                    "type": "string"
                },
                "terms_agree": {
                    "type": "boolean"
                }
            }
        },
        "models.UserLoginForm": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "use_cookie": {
                    "type": "boolean"
                }
            }
        }
    },
    "securityDefinitions": {
        "BEARER": {
            "description": "Type \"BEARER\" followed by a space and JWT token.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "app server",
	Description:      "app server",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

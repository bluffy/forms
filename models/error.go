package models

type AppError struct {
	Error struct {
		Fields  *map[string]string `json:"fields,omitempty"`
		Message *string            `json:"message,omitempty"`
		Code    *int               `json:"code,omitempty"`
	} `json:"error"`
}

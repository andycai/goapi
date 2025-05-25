package core

type Map struct {
	Key   string
	Value any
}

type List struct {
	Items []any
}

type Response struct {
	Code    int    `json:"code" default:"0"`
	Message string `json:"message" default:"success"`
	Data    any    `json:"data"`
}

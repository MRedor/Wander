package controllers

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func CreateError(code int, message string) Error {
	return Error{code, message}
}

package web

type webError struct {
	Code  int    `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}

var (
	InvalidParam = &webError{Code: -1400, Error: "invalid params"}
)

func NewWebError(code int, err string) *webError {
	return &webError{
		Code:  code,
		Error: err,
	}
}

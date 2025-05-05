package ex

import "encoding/json"

type ApiError struct {
	code    string
	message string
}

func New(code string, message string) error {
	return &ApiError{code, message}
}

func (e *ApiError) Error() string {
	return e.code + ": " + e.message
}

func (e *ApiError) MarshalJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *ApiError) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e)
}

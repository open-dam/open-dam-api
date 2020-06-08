package server

type ApiError struct {

	// A HTTP status code applicable to this problem
	Code float32 `json:"code"`

	// A description of the error that occurred
	Message string `json:"message"`
}

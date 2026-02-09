package i18n

import "errors"

var (
	Success             = "success"
	InternalServerError = "internal-server-error"
	BadRequest          = "bad-request"
	Unauthorized        = "unauthorized"
	Forbidden           = "forbidden"
	ValidateFailed      = "validate-failed"

	ExampleMessageOK = "example-message-ok"

	PrefixNotFound = "prefix-not-found"
	GenderNotFound = "gender-not-found"
)

var (
	ErrPrefixNotFound = errors.New(PrefixNotFound)
	ErrGenderNotFound = errors.New(GenderNotFound)
)

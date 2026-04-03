package errors

import "fmt"

const (
	ExitOK         = 0
	ExitGeneral    = 1
	ExitAuth       = 2
	ExitNotFound   = 3
	ExitValidation = 4
	ExitAPI        = 5
)

type ExitCoder interface {
	ExitCode() int
}

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("認証エラー: %s", e.Message)
}

func (e *AuthError) ExitCode() int {
	return ExitAuth
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%sが見つかりません: %s", e.Resource, e.ID)
}

func (e *NotFoundError) ExitCode() int {
	return ExitNotFound
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("入力エラー (%s): %s", e.Field, e.Message)
	}
	return fmt.Sprintf("入力エラー: %s", e.Message)
}

func (e *ValidationError) ExitCode() int {
	return ExitValidation
}

type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("APIエラー (HTTP %d, %s): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("APIエラー (HTTP %d): %s", e.StatusCode, e.Message)
}

func (e *APIError) ExitCode() int {
	return ExitAPI
}

type RateLimitError struct {
	Message string
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("APIレート制限: %s", e.Message)
}

func (e *RateLimitError) ExitCode() int {
	return ExitAPI
}

func GetExitCode(err error) int {
	if err == nil {
		return ExitOK
	}
	if ec, ok := err.(ExitCoder); ok {
		return ec.ExitCode()
	}
	return ExitGeneral
}

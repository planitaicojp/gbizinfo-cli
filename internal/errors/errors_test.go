package errors

import (
	"fmt"
	"testing"
)

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, ExitOK},
		{"auth", &AuthError{Message: "bad token"}, ExitAuth},
		{"not found", &NotFoundError{Resource: "法人", ID: "123"}, ExitNotFound},
		{"validation", &ValidationError{Field: "name", Message: "required"}, ExitValidation},
		{"api", &APIError{StatusCode: 500, Message: "server error"}, ExitAPI},
		{"rate limit", &RateLimitError{Message: "too many requests"}, ExitAPI},
		{"generic", fmt.Errorf("unknown"), ExitGeneral},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetExitCode(tt.err); got != tt.want {
				t.Errorf("GetExitCode() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"auth", &AuthError{Message: "invalid token"}, "認証エラー: invalid token"},
		{"not found", &NotFoundError{Resource: "法人", ID: "123"}, "法人が見つかりません: 123"},
		{"validation", &ValidationError{Field: "name", Message: "required"}, "入力エラー (name): required"},
		{"validation no field", &ValidationError{Message: "bad input"}, "入力エラー: bad input"},
		{"api", &APIError{StatusCode: 500, Message: "fail"}, "APIエラー (HTTP 500): fail"},
		{"api with code", &APIError{StatusCode: 400, Code: "BAD", Message: "fail"}, "APIエラー (HTTP 400, BAD): fail"},
		{"rate limit", &RateLimitError{Message: "limit exceeded"}, "APIレート制限: limit exceeded"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

package types

import (
	"errors"
)

// APIError API 错误类型
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error 实现 error 接口
func (e *APIError) Error() string {
	return e.Message
}

// NewBadRequestError 创建 400 错误
func NewBadRequestError(message string) error {
	return &APIError{
		Code:    400,
		Message: message,
	}
}

// NewUnauthorizedError 创建 401 错误
func NewUnauthorizedError(message string) error {
	return &APIError{
		Code:    401,
		Message: message,
	}
}

// NewForbiddenError 创建 403 错误
func NewForbiddenError(message string) error {
	return &APIError{
		Code:    403,
		Message: message,
	}
}

// NewNotFoundError 创建 404 错误
func NewNotFoundError(message string) error {
	return &APIError{
		Code:    404,
		Message: message,
	}
}

// NewInternalServerError 创建 500 错误
func NewInternalServerError(message string) error {
	return &APIError{
		Code:    500,
		Message: message,
	}
}

// GetErrorCode 获取错误码
func GetErrorCode(err error) int {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code
	}
	return 500 // 默认 500
}

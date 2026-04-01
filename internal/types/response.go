package types

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(data interface{}) *Response {
	return &Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(message string, data interface{}) *Response {
	return &Response{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// Error 错误响应
func Error(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

// BadRequest 400 错误
func BadRequest(message string) *Response {
	return &Response{
		Code:    400,
		Message: message,
		Data:    nil,
	}
}

// Unauthorized 401 错误
func Unauthorized(message string) *Response {
	return &Response{
		Code:    401,
		Message: message,
		Data:    nil,
	}
}

// Forbidden 403 错误
func Forbidden(message string) *Response {
	return &Response{
		Code:    403,
		Message: message,
		Data:    nil,
	}
}

// NotFound 404 错误
func NotFound(message string) *Response {
	return &Response{
		Code:    404,
		Message: message,
		Data:    nil,
	}
}

// InternalServerError 500 错误
func InternalServerError(message string) *Response {
	return &Response{
		Code:    500,
		Message: message,
		Data:    nil,
	}
}

// TooManyRequests 429 错误
func TooManyRequests(message string) *Response {
	return &Response{
		Code:    429,
		Message: message,
		Data:    nil,
	}
}

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Result 通用 REST 返回结构
type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// 常用业务码
const (
	CodeOk      = 0
	CodeFail    = 1
	CodeInvalid = 400
	CodeServer  = 500
)

// Ok 成功，带数据
func Ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Result{Code: CodeOk, Message: "success", Data: data})
}

// OkMsg 成功，自定义消息，无数据
func OkMsg(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Result{Code: CodeOk, Message: message})
}

// Fail 失败，业务码 + 消息 (根据业务码设置 HTTP 状态码)
func Fail(c *gin.Context, code int, message string) {
	httpStatus := http.StatusOK

	// 根据业务码映射 HTTP 状态码
	switch code {
	case 400:
		httpStatus = http.StatusBadRequest
	case 401:
		httpStatus = http.StatusUnauthorized
	case 403:
		httpStatus = http.StatusForbidden
	case 404:
		httpStatus = http.StatusNotFound
	case 500:
		httpStatus = http.StatusInternalServerError
	default:
		// 其他业务错误码使用 400
		if code >= 400 {
			httpStatus = code
		}
	}

	c.JSON(httpStatus, Result{Code: code, Message: message})
}

// FailBadRequest 参数错误 (HTTP 400)
func FailBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Result{Code: CodeInvalid, Message: message})
}

// FailServer 服务端错误 (HTTP 500)
func FailServer(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Result{Code: CodeServer, Message: message})
}

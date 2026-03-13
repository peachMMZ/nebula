package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func FromContext(c *gin.Context) *zap.Logger {

	id, exists := c.Get(RequestIDKey)

	if !exists {
		return Log
	}

	return Log.With(zap.String("request_id", id.(string)))
}

package http

import (
	"github.com/gin-gonic/gin"
)

func (h *HttpServer) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.log.Printf("HTTP %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

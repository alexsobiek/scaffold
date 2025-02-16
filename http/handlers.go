package http

import (
	"github.com/gin-gonic/gin"
)

func notFoundHandler(c *gin.Context) {
	NotFound(c, nil)
}

func methodNotAllowedHandler(c *gin.Context) {
	MethodNotAllowed(c, nil)
}

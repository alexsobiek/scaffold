package http

import "github.com/gin-gonic/gin"

func prepareKeys() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Keys == nil {
			c.Keys = make(map[string]interface{})
		}
		c.Next()
	}
}

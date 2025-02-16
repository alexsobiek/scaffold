package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func Ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Data: data})
}

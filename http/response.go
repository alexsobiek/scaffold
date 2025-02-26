package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Count int         `json:"count"`
	Page  int         `json:"page"`
}

func Ok[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, Response{Data: data})
}

func Paginated[T any](c *gin.Context, page int, data []T) {
	if len(data) == 0 {
		BadRequest(c, errors.New("no data"))
		return
	}

	count := len(data)
	c.JSON(http.StatusOK, PaginatedResponse{Data: data, Count: count, Page: page})
}

func Created[T any](c *gin.Context, data T) {
	c.JSON(http.StatusCreated, Response{Data: data})
}

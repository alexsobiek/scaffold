package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrInternal struct {
	Message string
}

func (e ErrInternal) Error() string {
	return e.Message
}

func InternalError(c *gin.Context, err error) {
	if err == nil || err.Error() == "" {
		err = ErrInternal{Message: "internal server error"}
	}

	c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
}

type ErrUnauthorized struct {
	Message string
}

func (e ErrUnauthorized) Error() string {
	return e.Message
}

func Unauthorized(c *gin.Context, err error) {
	if err == nil || err.Error() == "" {
		err = ErrUnauthorized{Message: "unauthorized"}
	}

	c.JSON(http.StatusUnauthorized, Response{Error: err.Error()})
}

type ErrNotFound struct {
	Message string
}

func (e ErrNotFound) Error() string {
	return e.Message
}

func NotFound(c *gin.Context, err error) {
	if err == nil || err.Error() == "" {
		err = ErrNotFound{Message: "not found"}
	}

	c.JSON(http.StatusNotFound, Response{Error: err.Error()})
}

type ErrMethodNotAllowed struct {
	Message string
}

func (e ErrMethodNotAllowed) Error() string {
	return e.Message
}

func MethodNotAllowed(c *gin.Context, err error) {
	if err == nil || err.Error() == "" {
		err = ErrMethodNotAllowed{Message: "method not allowed"}
	}

	c.JSON(http.StatusMethodNotAllowed, Response{Error: err.Error()})
}

type ErrBadRequest struct {
	Message string
}

func (e ErrBadRequest) Error() string {
	return e.Message
}

func BadRequest(c *gin.Context, err error) {
	if err == nil || err.Error() == "" {
		err = ErrBadRequest{Message: "bad request"}
	}

	c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
}

type ErrForbidden struct {
	Message string
}

func (e ErrForbidden) Error() string {
	return e.Message
}

func Forbidden(c *gin.Context, err error) {
	if err == nil || err.Error() == "" {
		err = ErrForbidden{Message: "forbidden"}
	}
	c.JSON(http.StatusForbidden, Response{Error: err.Error()})
}

func Error(c *gin.Context, err error) {
	switch err.(type) {
	case ErrInternal:
		InternalError(c, err)
	case ErrUnauthorized:
		Unauthorized(c, err)
	case ErrNotFound:
		NotFound(c, err)
	case ErrMethodNotAllowed:
		MethodNotAllowed(c, err)
	case ErrBadRequest:
		BadRequest(c, err)
	case ErrForbidden:
		Forbidden(c, err)
	default:
		InternalError(c, err)
	}
}

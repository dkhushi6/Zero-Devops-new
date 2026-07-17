package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

const RequestIDContextKey = "request_id"
const RequestIDHeader = "X-Request-Id"

func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		id := c.Request().Header.Get(RequestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}
		c.Set(RequestIDContextKey, id)
		c.Response().Header().Set(RequestIDHeader, id)
		return next(c)
	}
}

func GetRequestID(c *echo.Context) string {
	id, _ := c.Get(RequestIDContextKey).(string)
	return id
}

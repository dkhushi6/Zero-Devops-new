package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	echomw "github.com/labstack/echo/v5/middleware"
	"github.com/spf13/viper"
)

// NewCORS returns a CORS middleware configured with origins from the
// ALLOWED_ORIGINS environment variable.
func NewCORS() echo.MiddlewareFunc {
	allowedOrigins := strings.Split(
        viper.GetString("ALLOWED_ORIGINS"),
        ",",
    )

	return echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodDelete},
	})
}

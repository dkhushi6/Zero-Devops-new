package middleware

import (
	"Zero_Devops/server/domain"
	appmiddleware "Zero_Devops/server/middleware"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	UserIDContextKey = "user_id"
)

type AuthMiddlewareHandler struct {
	userRepo domain.UserRepository
}

func NewAuthMiddlewareHandler(userRepo domain.UserRepository) *AuthMiddlewareHandler {
	return &AuthMiddlewareHandler{
		userRepo: userRepo,
	}
}

func (a *AuthMiddlewareHandler) Validator(c *echo.Context, accessToken string) (int64, error) {
	secretKey := viper.GetString("JWT_SECRET")
	if secretKey == "" {
		return 0, domain.ErrMissingSecret
	}
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || token == nil || !token.Valid {
		return 0, domain.ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, domain.ErrInvalidToken
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, domain.ErrInvalidToken
	}
	userID := int64(userIDFloat)
	if a.userRepo != nil {
		existingUser, err := a.userRepo.GetByID(c.Request().Context(), userID)
		if err != nil {
			return 0, domain.ErrUserLookupFailed
		}
		if existingUser.ID == 0 {
			return 0, domain.ErrInvalidToken
		}
	}
	return userID, nil
}

func (a *AuthMiddlewareHandler) Skipper(c *echo.Context) bool {
	path := c.Path()
	return path == "/auth/github/login" || path == "/auth/refresh" || path == "/integration/scm/github/webhook"
}

func (a *AuthMiddlewareHandler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if a.Skipper(c) {
			return next(c)
		}
		
		log := appmiddleware.LoggerFromContext(c.Request().Context())

		cookie, err := c.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			log.Warn("Access token cookie missing")
			return echo.NewHTTPError(http.StatusUnauthorized, "access token cookie not found")
		}
		userID, err := a.Validator(c, cookie.Value)
		if err != nil {
			log.Warn("Token validation failed", zap.Error(err))
			switch {
			case errors.Is(err, domain.ErrMissingSecret):
				log.Error("JWT secret missing from configuration")
				return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
			case errors.Is(err, domain.ErrUserLookupFailed):
				log.Error("Database lookup for user failed during token validation")
				return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
			default:
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid access token")
			}
		}
		
		log.Info("User authenticated successfully", zap.Int64("user_id", userID))
		c.Set(UserIDContextKey, userID)
		return next(c)
	}
}

func (a *AuthMiddlewareHandler) ToMiddleware() echo.MiddlewareFunc {
	return a.AuthMiddleware
}

func GetUserID(c *echo.Context) (int64, bool) {
	userID, ok := c.Get(UserIDContextKey).(int64)
	return userID, ok
}

package http

import (
	"Zero_Devops/server/domain"
	"net/http"
	"github.com/sirupsen/logrus"
	"github.com/labstack/echo"
)

type ResponseError struct {
	Message string `json:"message"`
}

type AuthHandler struct {
	AUsecase domain.AuthUsecase
}

func NewAuthHandler(e *echo.Echo, us domain.AuthUsecase) {
	handler := &AuthHandler{
		AUsecase: us,
	}
	e.POST("/login", handler.Login)
	e.POST("/user/refresh",handler.Refresh)
	e.POST("/logout", handler.Logout)
	e.GET("/user/me", handler.GetUser)
}

func (a *AuthHandler) Login(c echo.Context) error {
	
	return nil
}

func (a* AuthHandler) Refresh(c echo.Context) error{
	return nil
}

func (a* AuthHandler) Logout(c echo.Context) error{
	return nil
}

func (a* AuthHandler) GetUser(c echo.Context) error{
	return nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

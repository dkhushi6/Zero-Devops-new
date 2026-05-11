package http

import (
	"Zero_Devops/server/domain"

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
	e.GET("/user/:id", handler.GetUser)
	e.POST("/logout", handler.Logout)
}

func (a *AuthHandler) Login(c echo.Context) error {

}

package handler

import (
	"splitz/config"
	"splitz/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthServiceInterface
	cnf     *config.EnvModel
}

func InitAuthHandler(service service.AuthServiceInterface, cnf *config.EnvModel) *AuthHandler {
	return &AuthHandler{service: service, cnf: cnf}
}

func (a *AuthHandler) Login(c *gin.Context) {
	a.service.Login(c)
}

func (a *AuthHandler) Register(c *gin.Context) {
	a.service.Register(c)
}

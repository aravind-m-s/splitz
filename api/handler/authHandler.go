package handler

import (
	"splitz/config"
	"splitz/service"

	"github.com/gin-gonic/gin"
)

type AuthHandlerStruct struct {
	service service.AuthServiceInterface
	cnf     *config.EnvModel
}

func AuthHandler(service service.AuthServiceInterface, cnf *config.EnvModel) *AuthHandlerStruct {
	return &AuthHandlerStruct{service: service, cnf: cnf}
}

func (a *AuthHandlerStruct) Login(c *gin.Context) {
	a.service.Login(c)
}

func (a *AuthHandlerStruct) Register(c *gin.Context) {
	a.service.Register(c)
}

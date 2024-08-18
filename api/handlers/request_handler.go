package handler

import (
	"splitz/common"
	"splitz/service"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	service service.RequestServiceInterface
	cnf     *common.JWTStruct
}

func InitRequestHandler(service service.RequestServiceInterface, cnf *common.JWTStruct) *RequestHandler {
	return &RequestHandler{service: service, cnf: cnf}
}

func (a *RequestHandler) CreateRequest(c *gin.Context) {
	a.service.CreateRequest(c, a.cnf)
}

func (a *RequestHandler) PayShare(c *gin.Context) {
	a.service.PayShare(c)
}

func (a *RequestHandler) ListRequest(c *gin.Context) {
	a.service.ListRequest(c)
}

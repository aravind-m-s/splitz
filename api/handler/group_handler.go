package handler

import (
	"net/http"
	"splitz/common"
	"splitz/service"

	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	service service.GroupServiceInterface
	cnf     *common.JWTStruct
}

func InitGroupHandler(service service.GroupServiceInterface, cnf *common.JWTStruct) *GroupHandler {
	return &GroupHandler{service: service, cnf: cnf}
}

func (a *GroupHandler) CreateGroup(c *gin.Context) {

	a.service.CreateGroup(c, a.cnf)
}

func (a *GroupHandler) DeleteGroup(c *gin.Context) {

	groupId := c.Param("id")

	a.service.DeleteGroup(groupId)
}

func (a *GroupHandler) GroupDetails(c *gin.Context) {

	groupId := c.Param("id")

	a.service.GroupDetails(groupId)
}

func (a *GroupHandler) ListGroup(c *gin.Context) {

	userId := c.Query("user_id")

	groups, err := a.service.ListGroup(userId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &groups)
}

func (a *GroupHandler) UpdateGroup(c *gin.Context) {

	groupId := c.Param("id")

	a.service.UpdateGroup(groupId)
}

package handler

import (
	"encoding/json"
	"net/http"
	"splitz/common"
	"splitz/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	s := c.Request.Header.Get("Authorization")

	token := strings.TrimPrefix(s, "Bearer ")
	user, err := a.cnf.GetFromToken(token, "user_id")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Unable to get user from token"})
		return
	}

	userId, uuidError := uuid.Parse(user)
	if uuidError != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Unable to parse user from token"})
		return
	}

	group, err := a.service.GroupDetails(groupId, userId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid group Id"})
		return
	}

	c.JSON(http.StatusOK, group)
}

func (a *GroupHandler) ListGroup(c *gin.Context) {

	userId := c.Query("user_id")

	if len(userId) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "user_id is required"})
		return
	}

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

func (a *GroupHandler) GetUserList(c *gin.Context) {

	contactsStr := c.PostForm("contacts")

	var contacts []string
	if err := json.Unmarshal([]byte(contactsStr), &contacts); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid contacts"})
		return
	}

	response, err := a.service.GetUserList(contacts)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (a *GroupHandler) CreateRequest(c *gin.Context) {

}

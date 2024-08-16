package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"splitz/common"
	"splitz/domain"
	"splitz/repository"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestServiceInterface interface {
	CreateRequest(c *gin.Context, cnf *common.JWTStruct) (data domain.GroupListResponse)
}

type requestServiceStruct struct {
	repo    repository.RequestInterface
	service GroupServiceInterface
}

func InitRequestService(repo repository.RequestInterface, service GroupServiceInterface) RequestServiceInterface {
	return &requestServiceStruct{repo: repo, service: service}
}

func (d *requestServiceStruct) CreateRequest(c *gin.Context, cnf *common.JWTStruct) (data domain.GroupListResponse) {

	s := c.Request.Header.Get("Authorization")

	token := strings.TrimPrefix(s, "Bearer ")

	owner, err := cnf.GetFromToken(token, "user_id")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unable to retrive user data"})
		return
	}

	requestType := c.PostForm("type")
	note := c.PostForm("note")
	amount := c.PostForm("amount")
	group := c.PostForm("group")

	fmt.Printf("uuid.UUID: %v\n", uuid.Max)

	groupId, groupErr := uuid.Parse(group)

	if groupErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid Group Id"})
		return

	}

	ownerId, ownerErr := uuid.Parse(owner)

	if ownerErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid Owner Id"})
		return
	}

	_, groupDetailsErr := d.service.GroupDetails(groupId.String(), ownerId)

	if groupDetailsErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Group doesn't exist"})
		return
	}

	usersJSON := c.PostForm("users")

	var users []map[string]string
	if err := json.Unmarshal([]byte(usersJSON), &users); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid users"})
		return
	}

	requestId, err := d.repo.CreateRequest(requestType, note, amount, groupId, ownerId)

	if err != nil {

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	for _, user := range users {
		userErr := d.repo.CreateUserRequest(requestId, user["amount"], user["id"])
		if userErr != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Unable to create user for request"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Request created successfully"})
	return
}

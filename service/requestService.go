package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"splitz/common"
	"splitz/repository"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestServiceInterface interface {
	CreateRequest(c *gin.Context, cnf *common.JWTStruct)
	PayShare(c *gin.Context)
	ListRequest(c *gin.Context)
}

type requestServiceStruct struct {
	repo    repository.RequestInterface
	service GroupServiceInterface
}

func InitRequestService(repo repository.RequestInterface, service GroupServiceInterface) RequestServiceInterface {
	return &requestServiceStruct{repo: repo, service: service}
}

func (d *requestServiceStruct) CreateRequest(c *gin.Context, cnf *common.JWTStruct) {

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

	createErr := d.repo.CreateRequest(requestType, note, amount, groupId, ownerId, users)

	if createErr != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, createErr)
	} else {
		c.JSON(http.StatusCreated, gin.H{"message": "Request created successfully"})
	}

	return
}

func (d *requestServiceStruct) PayShare(c *gin.Context) {

	request := c.Param("id")

	group := c.PostForm("group_id")
	user := c.PostForm("user_id")

	requestId, reqErr := uuid.Parse(request)

	if reqErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid Request Id"})
		return

	}

	groupId, groupErr := uuid.Parse(group)

	if groupErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid Group Id"})
		return

	}

	userId, userErr := uuid.Parse(user)

	if userErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid User Id"})
		return
	}

	amount := c.PostForm("amount")

	paidAmount, strConvErr := strconv.ParseFloat(amount, 64)

	if strConvErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid Amount"})
		return
	}

	err := d.repo.PayShare(requestId, groupId, userId, paidAmount)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Share paid successfully"})

	return
}

func (d *requestServiceStruct) ListRequest(c *gin.Context) {

	group := c.Param("id")

	groupId, groupErr := uuid.Parse(group)

	if groupErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid Group Id"})
		return

	}

	err, userReqs := d.repo.ListRequest(groupId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userReqs})

	return
}

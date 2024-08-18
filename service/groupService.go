package service

import (
	"encoding/json"
	"net/http"
	"splitz/common"
	"splitz/domain"
	"splitz/repository"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupServiceInterface interface {
	CreateGroup(c *gin.Context, cnf *common.JWTStruct) (data domain.GroupListResponse)
	DeleteGroup(id string) (data domain.GroupListResponse)
	GroupDetails(id string, userId uuid.UUID) (data domain.GroupDetailsResponse, errorMessage error)
	ListGroup(userId string) (data []domain.GroupListResponse, groupError error)
	UpdateGroup(id string) (data domain.GroupListResponse)
	GetUserList(contacts []string) (response gin.H, listError error)
}

type groupServiceStruct struct {
	repo repository.GroupRepoInterface
}

func InitGroupService(repo repository.GroupRepoInterface) GroupServiceInterface {
	return &groupServiceStruct{repo: repo}
}

func (d *groupServiceStruct) CreateGroup(c *gin.Context, cnf *common.JWTStruct) (data domain.GroupListResponse) {
	s := c.Request.Header.Get("Authorization")

	token := strings.TrimPrefix(s, "Bearer ")

	admin, err := cnf.GetFromToken(token, "user_id")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unable to retrive user data"})
		return
	}

	name := c.PostForm("name")
	imageData, _ := c.FormFile("image")
	usersStr := c.PostForm("users")

	if name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Name is required"})
		return
	} else if usersStr == "" || usersStr == "[]" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Atleast 1 user is required"})
		return
	}

	filePath := ""

	var users []string
	if err := json.Unmarshal([]byte(usersStr), &users); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid users"})
		return
	}

	for _, user := range users {
		status, err := d.repo.GetSingleUsers(user)
		if !status || err != "" {
			if err == "Internal Server Error" {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
				return

			} else {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User does not exist"})
				return

			}
		}
	}

	if imageData != nil {
		filePath = "./media/" + name + imageData.Filename
		if err := c.SaveUploadedFile(imageData, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

	}

	filePath = c.Request.Host + "/media/" + name + imageData.Filename

	createErr := d.repo.CreateGroup(name, filePath, admin, users)

	if createErr != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Unable to create group"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Group created successfully"})

	return
}

func (d *groupServiceStruct) DeleteGroup(id string) (data domain.GroupListResponse) {
	return d.repo.DeleteGroup(id)
}

func (d *groupServiceStruct) GroupDetails(id string, userId uuid.UUID) (data domain.GroupDetailsResponse, errorMessage error) {

	return d.repo.GroupDetails(id, userId)
}

func (d *groupServiceStruct) ListGroup(userId string) (data []domain.GroupListResponse, groupError error) {

	return d.repo.ListGroup(userId)
}

func (d *groupServiceStruct) UpdateGroup(id string) (data domain.GroupListResponse) {
	return d.repo.UpdateGroup(id)
}

func (d *groupServiceStruct) GetUserList(contacts []string) (response gin.H, listError error) {

	response, listError = d.repo.GetUserList(contacts)

	return response, listError

}

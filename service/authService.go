package service

import (
	"net/http"
	"splitz/common"
	"splitz/domain"
	"splitz/repository"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type AuthServiceInterface interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type authServiceStruct struct {
	repo repository.AuthRepoInterface
	jwt  *common.JWTStruct
}

func InitAuthService(repo repository.AuthRepoInterface, jwt *common.JWTStruct) AuthServiceInterface {
	return &authServiceStruct{repo: repo, jwt: jwt}
}

func (a *authServiceStruct) Login(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
		}
	}()

	mobile := c.PostForm("mobile")
	password := c.PostForm("password")

	if mobile == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Mobile is required"})
		return
	} else if password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Password is required"})
		return

	}

	user, err := a.repo.Login(mobile, password)

	if len(err) != 0 {
		statusCode := http.StatusBadRequest
		if err == "Internal Server error" {
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, gin.H{
			"message": err,
		})
	} else {
		responseUser := user.ToUserResponse()
		token, err := a.jwt.GenerateJWT(user.ID)

		if err != nil {
			println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to generate token",
			})
		} else {
			responseUser.Token = token
			c.JSON(http.StatusOK, &responseUser)
		}
	}
}

func (a *authServiceStruct) Register(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
		}
	}()

	file, _ := c.FormFile("image")
	mobile := c.PostForm("mobile")
	password := c.PostForm("password")
	name := c.PostForm("name")
	fcmTokens := c.PostForm("fcm_token")

	if mobile == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Mobile is required"})
		return
	} else if password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Password is required"})
		return
	} else if name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Name is required"})
		return
	} else if fcmTokens == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "FCM Token is required"})
		return
	}

	filePath := ""
	if file != nil {
		filePath = "./media/" + mobile + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	user := domain.User{
		Mobile:    mobile,
		Password:  password,
		Name:      name,
		Image:     c.Request.Host + "/media/" + mobile + file.Filename,
		FcmTokens: pq.StringArray{fcmTokens},
	}

	id, err := a.repo.Register(user)

	if len(err) != 0 {
		statusCode := http.StatusBadRequest
		if err == "Internal Server error" {
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, gin.H{
			"message": err,
		})
	} else {

		responseUser := user.ToUserResponse()
		token, err := a.jwt.GenerateJWT(id)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to generate token",
			})
		} else {
			responseUser.Token = token
			responseUser.ID = id
			c.JSON(http.StatusOK,
				&responseUser,
			)
		}

	}

}

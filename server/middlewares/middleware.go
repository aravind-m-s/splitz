package middlewares

import (
	"net/http"
	"splitz/common"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthorizationStruct struct {
	cnf *common.JWTStruct
}

func NewAuthorization(cnf *common.JWTStruct) *AuthorizationStruct {
	return &AuthorizationStruct{cnf: cnf}
}

func (cnf *AuthorizationStruct) AuthorizationMiddleware(c *gin.Context) {
	s := c.Request.Header.Get("Authorization")

	token := strings.TrimPrefix(s, "Bearer ")

	if len(token) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is required"})
		return
	}

	if _, err := cnf.cnf.VerifyJWT(token); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token expired please login again"})
		return
	}
}

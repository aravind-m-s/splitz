package api

import (
	"net/http"
	handler "splitz/api/handlers"
	"splitz/api/middlewares"
	"splitz/config"

	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewHandlerFunc(w http.ResponseWriter, r *http.Request) {

}

func Handler(authHandler *handler.AuthHandlerStruct, groupHandler *handler.GroupHandler, requestHandler *handler.RequestHandler, middleWare *middlewares.AuthorizationStruct) *ServerHTTP {
	engine := gin.New()

	engine.Use(gin.Logger())
	engine.Static("/media", "./media")

	// API Group
	apiGroup := engine.Group("/api")

	// Authorizaiton
	apiGroup.POST("/login", authHandler.Login)
	apiGroup.POST("/register", authHandler.Register)

	// Group
	groupEngine := apiGroup.Group("/group", middleWare.AuthorizationMiddleware)

	groupEngine.GET("/list", groupHandler.ListGroup)
	groupEngine.POST("/create", groupHandler.CreateGroup)
	groupEngine.GET("/:id", groupHandler.GroupDetails)
	groupEngine.PUT("/:id", groupHandler.UpdateGroup)
	groupEngine.DELETE("/:id", groupHandler.DeleteGroup)
	groupEngine.POST("/users/list", groupHandler.GetUserList)

	requestEngine := groupEngine.Group("/request")

	requestEngine.POST("/create", requestHandler.CreateRequest)
	requestEngine.POST("/pay/:id", requestHandler.PayShare)
	requestEngine.GET("/list/:id", requestHandler.ListRequest)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start(cnf *config.EnvModel) {
	sh.engine.Run(cnf.Port)
}

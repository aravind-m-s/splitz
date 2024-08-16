package di

import (
	"splitz/api"
	handler "splitz/api/handlers"
	"splitz/api/middlewares"
	"splitz/common"
	"splitz/config"
	"splitz/database"
	"splitz/repository"
	"splitz/service"
)

func InitServer(cnf *config.EnvModel) (*api.ServerHTTP, error) {
	db, err := database.InitDatabase(cnf)

	jwt := common.NewHelper(cnf)
	authorization := middlewares.NewAuthorization(jwt)

	// Authorization
	authRepo := repository.InitAuthRepo(db)
	authService := service.InitAuthService(authRepo, jwt)
	authHandler := handler.AuthHandler(authService, cnf)

	// Group
	groupRepo := repository.InitGroupRepo(db)
	groupService := service.InitGroupService(groupRepo)
	groupHandler := handler.InitGroupHandler(groupService, jwt)

	// Request
	requestRepo := repository.InitRequestRepo(db)
	requestService := service.InitRequestService(requestRepo, groupService)
	requestHandler := handler.InitRequestHandler(requestService, jwt)

	server := api.Handler(authHandler, groupHandler, requestHandler, authorization)

	if err != nil {
		return nil, err
	}

	return server, nil
}

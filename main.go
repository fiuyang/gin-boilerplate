package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"scylla/controller"
	"scylla/docs"
	"scylla/pkg/config"
	"scylla/pkg/exception"
	"scylla/pkg/helper"
	"scylla/pkg/utils"
	"scylla/repository"
	"scylla/router"
	"scylla/service"
	"time"
)

//	@title			Boilerplate API
//	@version		1.0
//	@description	Boilerplate API in Go using Gin framework

// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func main() {

	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	//Database
	db := config.ConnectionDB(&loadConfig)

	//Validate
	validate := utils.InitializeValidator(db)

	//Swagger
	if loadConfig.Environment != "dev" {
		docs.SwaggerInfo.Host = loadConfig.SwaggerHost
		docs.SwaggerInfo.BasePath = loadConfig.SwaggerUrl
	} else {
		docs.SwaggerInfo.Host = "localhost:8000"
		docs.SwaggerInfo.BasePath = "/api/v1"
	}

	//Init Repository
	userRepo := repository.NewUserRepoImpl(db)
	passResetRepo := repository.NewPassResetRepoImpl(db)
	customerRepo := repository.NewCustomerRepoImpl(db)

	//Init Service
	authService := service.NewAuthServiceImpl(userRepo, passResetRepo, validate)
	customerService := service.NewCustomerServiceImpl(customerRepo, validate)
	userSevice := service.NewUserServiceImpl(userRepo, validate)

	//Init controller
	authController := controller.NewAuthController(authService)
	customerController := controller.NewCustomerController(customerService)
	userController := controller.NewUserController(userSevice)

	//Router
	routes := router.NewRouter(
		authController,
		customerController,
		userController,
	)

	app := gin.Default()
	app.Use(gin.Logger())
	app.Use(gin.CustomRecovery(exception.ErrorHandlers))

	//docs swagger
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	app.Use(func(c *gin.Context) {
		routes.ServeHTTP(c.Writer, c.Request)
	})

	server := &http.Server{
		Addr:           ":" + loadConfig.ServerPort,
		Handler:        app,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	error := server.ListenAndServe()
	helper.ErrorPanic(error)
}

package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"scylla/controller"
	"scylla/docs"
	"scylla/dto"
	"scylla/pkg/config"
	"scylla/pkg/connection"
	"scylla/pkg/exception"
	"scylla/pkg/helper"
	"scylla/pkg/middleware"
	"scylla/pkg/utils"
	"scylla/repository"
	"scylla/service"
	"time"
)

//	@title			Boilerplate API
//	@version		1.0
//	@description	Boilerplate API in Go using Gin framework

// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description				   Type "Bearer" followed by a space and JWT token.
func main() {
	//config
	conf := config.Get()

	//database
	db := connection.GetDatabase(conf.Database)

	//validate
	validate := utils.InitializeValidator()

	//environment swagger
	if conf.Swagger.Mode != "dev" {
		docs.SwaggerInfo.Host = conf.Swagger.Host
		docs.SwaggerInfo.BasePath = conf.Swagger.Url
	} else {
		docs.SwaggerInfo.Host = "localhost:8000"
		docs.SwaggerInfo.BasePath = "/api/v1"
	}

	//repository
	userRepo := repository.NewUserRepoImpl(db)
	passResetRepo := repository.NewPassResetRepoImpl(db)
	customerRepo := repository.NewCustomerRepoImpl(db)

	//service
	authService := service.NewAuthServiceImpl(conf, userRepo, passResetRepo, validate)
	customerService := service.NewCustomerServiceImpl(customerRepo, validate)
	userService := service.NewUserServiceImpl(userRepo, validate)

	//controller
	authController := controller.NewAuthController(authService)
	customerController := controller.NewCustomerController(customerService)
	userController := controller.NewUserController(userService)

	//environment gin
	app := gin.Default()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())
	app.Use(gin.CustomRecovery(exception.ExceptionHandlers))
	app.Use(middleware.TracingMiddleware())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	app.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dto.Error{
			Code:    http.StatusNotFound,
			Status:  "NOT FOUND",
			Errors:  "Page Not Found",
			TraceID: uuid.New().String(),
		})
	})
	//docs swagger
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//routes v1
	authController.Route(app)
	customerController.Route(app)
	userController.Route(app)

	//server
	server := &http.Server{
		Addr:           ":" + conf.Server.Port,
		Handler:        app,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	error := server.ListenAndServe()
	helper.ErrorPanic(error)
}

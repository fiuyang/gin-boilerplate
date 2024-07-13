package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"scylla/controller"
	"scylla/pkg/middleware"
)

func NewRouter(
	authController *controller.AuthController,
	customerController *controller.CustomerController,
	userController *controller.UserController,
) *gin.Engine {

	app := gin.New()
	app.Use(middleware.TracingMiddleware())

	app.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":     http.StatusNotFound,
			"status":   "NOT FOUND",
			"errors":   "Page Not Found",
			"trace_id": uuid.New().String(),
		})
	})

	router := app.Group("/api/v1")

	//auth
	authRouter := router.Group("/auth")
	authRouter.POST("/register", authController.Register)
	authRouter.POST("/login", authController.Login)
	authRouter.POST("/forgot-password", authController.ForgotPassword)
	authRouter.POST("/check-otp", authController.CheckOtp)
	authRouter.PATCH("/reset-password", authController.ResetPassword)
	authRouter.POST("/logout", middleware.JwtMiddleware(), authController.Logout)

	router.Use(middleware.JwtMiddleware())
	//customer
	customerRouter := router.Group("/customers")
	customerRouter.GET("", customerController.FindAllPaging)
	customerRouter.GET("/:customerId", customerController.FindById)
	customerRouter.POST("", customerController.Create)
	customerRouter.POST("/batch", customerController.CreateBatch)
	customerRouter.PATCH("/:customerId", customerController.Update)
	customerRouter.DELETE("/batch", customerController.DeleteBatch)
	customerRouter.GET("/export", customerController.Export)
	customerRouter.POST("/import", customerController.Import)

	//user
	userRouter := router.Group("/users")
	userRouter.POST("", userController.Create)
	userRouter.PATCH("/:userId", userController.Update)
	userRouter.GET("/:userId", userController.FindById)
	userRouter.GET("", userController.FindAll)
	userRouter.POST("/batch", userController.DeleteBatch)
	userRouter.GET("/export", userController.Export)
	userRouter.POST("/import", userController.Import)

	return app
}

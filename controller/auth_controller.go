package controller

import (
	"context"
	"net/http"
	"scylla/dto"
	"scylla/pkg/exception"
	"scylla/pkg/helper"
	"scylla/pkg/middleware"
	"scylla/pkg/utils"
	"scylla/service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (controller *AuthController) Route(app *gin.Engine) {
	authRouter := app.Group("/api/v1/auth")
	authRouter.POST("/register", controller.Register)
	authRouter.POST("/login", controller.Login)
	authRouter.POST("/forgot-password", controller.ForgotPassword)
	authRouter.POST("/check-otp", controller.CheckOtp)
	authRouter.PATCH("/reset-password", controller.ResetPassword)
	authRouter.POST("/logout", middleware.JwtMiddleware(), controller.Logout)
}

// Note		godoc
//
// @Summary		Login
// @Description	Login.
// @Param		data	body	dto.LoginRequest	true	"login"
// @Produce		application/json
// @Tags		auth
// @Success		200	{object}	dto.JsonSuccess{data=string}	"Data"
// @Failure		400	{object}	dto.JsonBadRequest{}			"Validation error"
// @Failure		404	{object}	dto.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	dto.JsonInternalServerError{}	"Internal server error"
// @Router		/auth/login [post]
func (controller *AuthController) Login(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := dto.LoginRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	token, err := controller.authService.Login(c, request)

	helper.ErrorPanic(err)

	webResponse := dto.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "Login Successful",
		Data:    token,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// Note		godoc
//
// @Summary		Register
// @Description	Register.
// @Param		data	body	dto.CreateUserRequest	true	"register"
// @Produce		application/json
// @Tags		auth
// @Success		201	{object}	dto.JsonCreated{data=nil}		"Data"
// @Failure		400	{object}	dto.JsonBadRequest{}			"Validation error"
// @Failure		404	{object}	dto.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	dto.JsonInternalServerError{}	"Internal server error"
// @Router		/auth/register [post]
func (controller *AuthController) Register(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := dto.CreateUserRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	controller.authService.Register(c, request)

	webResponse := dto.Response{
		Code:    http.StatusCreated,
		Status:  "Ok",
		Message: "Create User Successful",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusCreated, webResponse)
}

// Note		godoc
//
// @Summary		Forgot Password
// @Description	Forgot Password.
// @Param		data	body	dto.ForgotPasswordRequest	true	"forgot password"
// @Produce		application/json
// @Tags		auth
// @Success		200	{object}	dto.JsonSuccess{data=string}	"Data"
// @Failure		400	{object}	dto.JsonBadRequest{}			"Validation error"
// @Failure		404	{object}	dto.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	dto.JsonInternalServerError{}	"Internal server error"
// @Router		/auth/forgot-password [post]
func (controller *AuthController) ForgotPassword(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := dto.ForgotPasswordRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	otp, err := controller.authService.ForgotPassword(c, request)
	helper.ErrorPanic(err)

	webResponse := dto.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Forgot Password Successful",
		Data:    otp,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// Note		godoc
//
// @Summary		Check Otp
// @Description	Check Otp.
// @Param		data	body	dto.CheckOtpRequest	true	"check otp"
// @Produce		application/json
// @Tags		auth
// @Success		200	{object}	dto.JsonSuccess{data=nil}		"Data"
// @Failure		400	{object}	dto.JsonBadRequest{}			"Validation error"
// @Failure		404	{object}	dto.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	dto.JsonInternalServerError{}	"Internal server error"
// @Router		/auth/check-otp [post]
func (controller *AuthController) CheckOtp(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := dto.CheckOtpRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	message, err := controller.authService.CheckOtp(c, request)
	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}

	webResponse := dto.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: message,
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// Note		godoc
//
// @Summary		Reset Password
// @Description	Reset Password.
// @Param		data	body	dto.ResetPasswordRequest	true	"reset password"
// @Produce		application/json
// @Tags		auth
// @Success		200	{object}	dto.JsonSuccess{data=nil}		"Data"
// @Failure		400	{object}	dto.JsonBadRequest{}			"Validation error"
// @Failure		404	{object}	dto.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	dto.JsonInternalServerError{}	"Internal server error"
// @Router		/auth/reset-password [patch]
func (controller *AuthController) ResetPassword(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := dto.ResetPasswordRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	message, err := controller.authService.ResetPassword(c, request)
	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}

	webResponse := dto.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: message,
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// Note		godoc
//
// @Summary		Logout
// @Description	Logout.
// @Produce		application/json
// @Tags		auth
// @Success		200	{object}	dto.JsonSuccess{data=nil}		"Data"
// @Failure		400	{object}	dto.JsonBadRequest{}			"Validation error"
// @Failure		404	{object}	dto.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	dto.JsonInternalServerError{}	"Internal server error"
// @Router		/auth/logout [post]
// @Security	Bearer
func (controller *AuthController) Logout(ctx *gin.Context) {
	token := extractTokenFromRequest(ctx)

	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty token"})
		return
	}

	err := controller.authService.Logout(token)
	helper.ErrorPanic(err)

	webResponse := dto.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Logout Successful",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

func extractTokenFromRequest(ctx *gin.Context) string {
	token := ctx.GetHeader("Authorization")
	if token != "" {
		return strings.TrimPrefix(token, "Bearer ")
	}
	return ""
}

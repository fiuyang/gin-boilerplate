package controller

import (
	"context"
	"net/http"
	"scylla/entity"
	"scylla/pkg/helper"
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

// CreateTags		godoc
//
// @Summary		Login
// @Description	Login.
// @Param			data	body	entity.LoginRequest	true	"login"
// @Produce		application/json
// @Tags			auth
// @Success		200	{object}	entity.JsonSuccess{data=string}		"Data"
// @Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
// @Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
// @Router			/auth/login [post]
func (controller *AuthController) Login(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := entity.LoginRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	token, err := controller.authService.Login(c, request)

	helper.ErrorPanic(err)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "Login Successful",
		Data:    token,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
//
// @Summary		Register
// @Description	Register.
// @Param			data	body	entity.CreateUserRequest	true	"register"
// @Produce		application/json
// @Tags			auth
// @Success		201	{object}	entity.JsonCreated{data=nil}		"Data"
// @Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
// @Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
// @Router			/auth/register [post]
func (controller *AuthController) Register(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := entity.CreateUserRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	controller.authService.Register(c, request)

	webResponse := entity.Response{
		Code:    http.StatusCreated,
		Status:  "Ok",
		Message: "Create User Successful",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusCreated, webResponse)
}

// CreateTags		godoc
//
// @Summary		Forgot Password
// @Description	Forgot Password.
// @Param			data	body	entity.ForgotPasswordRequest	true	"forgot password"
// @Produce		application/json
// @Tags			auth
// @Success		200	{object}	entity.JsonSuccess{data=string}		"Data"
// @Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
// @Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
// @Router			/auth/forgot-password [post]
func (controller *AuthController) ForgotPassword(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := entity.ForgotPasswordRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	otp, err := controller.authService.ForgotPassword(c, request)
	helper.ErrorPanic(err)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Forgot Password Successful",
		Data:    otp,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
//
// @Summary		Check Otp
// @Description	Check Otp.
// @Param			data	body	entity.CheckOtpRequest	true	"check otp"
// @Produce		application/json
// @Tags			auth
// @Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
// @Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
// @Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
// @Router			/auth/check-otp [post]
func (controller *AuthController) CheckOtp(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := entity.CheckOtpRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	controller.authService.CheckOtp(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Otp Valid",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
//
// @Summary		Reset Password
// @Description	Reset Password.
// @Param			data	body	entity.ResetPasswordRequest	true	"reset password"
// @Produce		application/json
// @Tags			auth
// @Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
// @Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
// @Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
// @Router			/auth/reset-password [patch]
func (controller *AuthController) ResetPassword(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := entity.ResetPasswordRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	controller.authService.ResetPassword(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Reset Password Successful",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
//
// @Summary		Logout
// @Description	Logout.
// @Produce		application/json
// @Tags			auth
// @Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
// @Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
// @Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
// @Router			/auth/logout [post]
// @Security		Bearer
func (controller *AuthController) Logout(ctx *gin.Context) {
	token := extractTokenFromRequest(ctx)

	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty token"})
		return
	}

	err := controller.authService.Logout(token)
	helper.ErrorPanic(err)

	webResponse := entity.Response{
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

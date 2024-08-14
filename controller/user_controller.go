package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"scylla/dto"
	"scylla/pkg/exception"
	"scylla/pkg/helper"
	"scylla/pkg/middleware"
	"scylla/pkg/utils"
	"scylla/service"
	"time"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (controller *UserController) Route(app *gin.Engine) {
	userRouter := app.Group("/api/v1/users", middleware.JwtMiddleware())
	userRouter.POST("", controller.Create)
	userRouter.PATCH("/:userId", controller.Update)
	userRouter.GET("/:userId", controller.FindById)
	userRouter.GET("", controller.FindAll)
	userRouter.POST("/batch", controller.DeleteBatch)
	userRouter.GET("/export", controller.Export)
	userRouter.POST("/import", controller.Import)
}

// Note		godoc
//
//	@Summary		Create user
//	@Description	Create user.
//	@Param			data	body	dto.CreateUserRequest	true	"create user"
//	@Produce		application/json
//	@Tags			users
//	@Success		201	{object}	dto.JsonCreated{data=nil}       "Data"
//	@Failure		400	{object}	dto.JsonBadRequest{}			"Validation error"
//	@Failure		404	{object}	dto.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	dto.JsonInternalServerError{}	"Internal server error"
//	@Router			/users [post]
//	@Security		Bearer
func (controller *UserController) Create(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := dto.CreateUserRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	controller.userService.Create(c, request)

	webResponse := dto.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Create Successful",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// Note		godoc
//
//	@Summary		Update user
//	@Description	Update user.
//	@Param			userId	path	string					true	"user_id"
//	@Param			data	body	dto.UpdateUserRequest	true	"update user"
//	@Tags			users
//	@Produce		application/json
//	@Success		200	{object}	dto.JsonSuccess{data=nil}       "Data"
//	@Failure		400	{object}	dto.JsonBadRequest{}			"Validation error"
//	@Failure		404	{object}	dto.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	dto.JsonInternalServerError{}	"Internal server error"
//	@Router			/users/{userId} [patch]
//	@Security		Bearer
func (controller *UserController) Update(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := dto.UpdateUserRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	var params dto.UserParams

	if err := ctx.ShouldBindUri(&params); err != nil {
		panic(exception.NewBadRequestHandler(err.Error()))
	}

	request.ID = params.UserId

	controller.userService.Update(c, request)

	webResponse := dto.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "Update Successful",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// Note		godoc
//
//	@Summary		Delete batch user
//	@Description	Delete batch user.
//	@Param			data	body	dto.DeleteBatchUserRequest	true	"delete batch user"
//	@Produce		application/json
//	@Tags			users
//	@Success		200	{object}	dto.JsonSuccess{data=nil}       "Data"
//	@Failure		400	{object}	dto.JsonBadRequest{}			"Validation error"
//	@Failure		404	{object}	dto.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	dto.JsonInternalServerError{}	"Internal server error"
//	@Router			/users/batch [post]
//	@Security		Bearer
func (controller *UserController) DeleteBatch(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := dto.DeleteBatchUserRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	controller.userService.DeleteBatch(c, request)

	webResponse := dto.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "Delete Batch Successful",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// Note 		godoc
//
//	@Summary		Get user by id.
//	@Description	Get user by id.
//	@Param			userId	path	string	true	"user_id"
//	@Produce		application/json
//	@Tags			users
//	@Success		200	{object}	dto.JsonSuccess{data=nil}       "Data"
//	@Failure		400	{object}	dto.JsonBadRequest{}			"Validation error"
//	@Failure		404	{object}	dto.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	dto.JsonInternalServerError{}	"Internal server error"
//	@Router			/users/{userId} [get]
//	@Security		Bearer
func (controller *UserController) FindById(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var params dto.UserParams

	if err := ctx.ShouldBindUri(&params); err != nil {
		panic(exception.NewBadRequestHandler(err.Error()))
	}

	response := controller.userService.FindById(c, params)

	webResponse := dto.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   response,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// Note		godoc
//
//	@Summary		Get all users.
//	@Description	Get all users.
//	@Produce		application/json
//	@Tags			users
//	@Param			start_date	query		string		false	"start_date"
//	@Param			end_date	query		string		false	"end_date"
//	@Param			username	query		string	    false	"username"
//	@Param			email		query		string	    false	"email"
//	@Success		200			{object}	dto.Response{data=[]dto.UserResponse{}}   	"Data"
//	@Failure		400			{object}	dto.JsonBadRequest{}						"Validation error"
//	@Failure		404			{object}	dto.JsonNotFound{}							"Data not found"
//	@Failure		500			{object}	dto.JsonInternalServerError{}				"Internal server error"
//	@Router			/users [get]
//	@Security		Bearer
func (controller *UserController) FindAll(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var dataFilter dto.UserQueryFilter

	if err := ctx.ShouldBindQuery(&dataFilter); err != nil {
		panic(exception.NewBadRequestHandler(err.Error()))
	}

	response := controller.userService.FindAll(c, dataFilter)

	webResponse := dto.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// Note 		    godoc
//
//	@Summary		Export Excel user.
//	@Description	Export Excel user.
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Tags			users
//	@Param			start_date	query		string	false	"start_date"
//	@Param			end_date	query		string	false	"end_date"
//	@Param			username	query		string	false	"username"
//	@Param			email		query		string	false	"email"
//	@Success		200			{object}	dto.JsonSuccess{data=string}    "Data"
//	@Failure		400			{object}	dto.JsonBadRequest{}			"Validation error"
//	@Failure		404			{object}	dto.JsonNotFound{}				"Data not found"
//	@Failure		500			{object}	dto.JsonInternalServerError{}	"Internal server error"
//	@Router			/users/export [get]
//	@Security		Bearer
func (controller *UserController) Export(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var dataFilter dto.UserQueryFilter

	if err := ctx.ShouldBindQuery(&dataFilter); err != nil {
		panic(exception.NewBadRequestHandler(err.Error()))
	}

	filePath, err := controller.userService.Export(c, dataFilter)
	helper.ErrorPanic(err)
	defer os.Remove(filePath) // Remove the file after the function exits

	fileName := filepath.Base(filePath)
	// Set headers for the Excel file
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	// Read the Excel file and write to the response body
	data, err := os.ReadFile(filePath)
	helper.ErrorPanic(err)

	// Write data to the response body
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

//	Note 		    godoc
//
// @Summary		Import Excel user.
// @Description	Import Excel user.
// @Produce		application/json
// @Tags		users
// @Param		data	formData	file	true	"import Excel user"
// @Success		200		{object}	dto.JsonSuccess{data=string}    "Data"
// @Failure		400		{object}	dto.JsonBadRequest{}			"Validation error"
// @Failure		404		{object}	dto.JsonNotFound{}				"Data not found"
// @Failure		500		{object}	dto.JsonInternalServerError{}	"Internal server error"
// @Router		/users/import [post]
// @Security	Bearer
func (controller *UserController) Import(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	file, err := ctx.FormFile("file")
	if err != nil {
		panic(exception.NewBadRequestHandler(err.Error()))
	}

	// Call the usersService method to handle the import
	err = controller.userService.Import(c, file)
	helper.ErrorPanic(err)

	webResponse := dto.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "Import Successful",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

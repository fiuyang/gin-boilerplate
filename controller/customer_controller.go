package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"scylla/entity"
	"scylla/pkg/exception"
	"scylla/pkg/helper"
	"scylla/pkg/utils"
	"scylla/service"
	"time"
)

type CustomerController struct {
	customerService service.CustomerService
}

func NewCustomerController(customerService service.CustomerService) *CustomerController {
	return &CustomerController{
		customerService: customerService,
	}
}

// CreateTags		godoc
//
//	@Summary		Create customer
//	@Description	Create customer.
//	@Param			data	formData	entity.CreateCustomerRequest	true	"create customer"
//	@Produce		application/json
//	@Tags			customers
//	@Success		201	{object}	entity.JsonCreated{data=nil}"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/customers [post]
//	@Security		Bearer
func (handler *CustomerController) Create(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := entity.CreateCustomerRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	err = handler.customerService.Create(c, request)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	webResponse := entity.Response{
		Code:    http.StatusCreated,
		Status:  "Created",
		Message: "Created Successful",
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
//
//	@Summary		Create customer batch
//	@Description	Create customer batch.
//	@Param			data	body	entity.CreateCustomerBatchRequest	true	"create customer batch"
//	@Produce		application/json
//	@Tags			customers
//	@Success		201	{object}	entity.JsonCreated{data=nil}"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/customers/batch [post]
//	@Security		Bearer
func (handler *CustomerController) CreateBatch(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := entity.CreateCustomerBatchRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	handler.customerService.CreateBatch(c, request)

	webResponse := entity.Response{
		Code:    http.StatusCreated,
		Status:  "Created",
		Message: "Created Batch Successful",
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
//
//	@Summary		update customer
//	@Description	update customer.
//	@Param			data		body	entity.UpdateCustomerRequest	true	"update customer"
//	@Param			customerId	path	string							true	"customer_id"
//	@Produce		application/json
//	@Tags			customers
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/customers/{customerId} [patch]
//	@Security		Bearer
func (handler *CustomerController) Update(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := entity.UpdateCustomerRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	var params entity.CustomerParams

	if err := ctx.ShouldBindUri(&params); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	request.ID = params.CustomerId

	handler.customerService.Update(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Update Successful",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
//
//	@Summary		Delete batch customer
//	@Description	Delete batch customer.
//	@Param			data	body	entity.DeleteBatchCustomerRequest	true	"delete batch customer"
//	@Produce		application/json
//	@Tags			customers
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/customers/batch [delete]
//	@Security		Bearer
func (handler *CustomerController) DeleteBatch(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := entity.DeleteBatchCustomerRequest{}
	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	handler.customerService.DeleteBatch(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Delete Batch Successful",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
//
//	@Summary		get customer by id.
//	@Param			customerId	path	string	true	"customer_id"
//	@Description	get customer by id.
//	@Produce		application/json
//	@Tags			customers
//	@Success		200	{object}	entity.JsonSuccess{data=entity.CustomerResponse{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/customers/{customerId} [get]
//	@Security		Bearer
func (handler *CustomerController) FindById(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var params entity.CustomerParams

	if err := ctx.ShouldBindUri(&params); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	data := handler.customerService.FindById(c, params)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
//
//	@Summary		Get all customers.
//	@Description	Get all customers.
//	@Produce		application/json
//	@Param			limit		query	string	false	"limit"
//	@Param			page		query	string	false	"page"
//	@Param			start_date	query	string	false	"start_date"
//	@Param			username	query	string	false	"username"
//	@Param			email		query	string	false	"email"
//	@Param			end_date	query	string	false	"end_date"
//	@Param			sort		query	string	false	"sort"
//	@Tags			customers
//	@Success		200	{object}	entity.Response{data=[]entity.CustomerResponse{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/customers [get]
//	@Security		Bearer
func (handler *CustomerController) FindAllPaging(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var dataFilter entity.CustomerQueryFilter

	if err := ctx.ShouldBindQuery(&dataFilter); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	response, paging := handler.customerService.FindAllPaging(c, dataFilter)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
		Meta:   &paging,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// ExportTags 		    godoc
//
//	@Summary		Export Excel customer.
//	@Description	Export Excel customer.
//	@Produce		application/json
//	@Tags			customers
//	@Param			start_date	query		string	false	"start_date"
//	@Param			end_date	query		string	false	"end_date"
//	@Param			username	query		string	false	"username"
//	@Param			email		query		string	false	"email"
//	@Success		200			{object}	entity.JsonSuccess{data=string}"Data"
//	@Failure		400			{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404			{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500			{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/customers/export [get]
//	@Security		Bearer
func (controller *CustomerController) Export(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var dataFilter entity.CustomerQueryFilter

	if err := ctx.ShouldBindQuery(&dataFilter); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	filePath, err := controller.customerService.Export(c, dataFilter)
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

// ImportTags 		    godoc
//
//	@Summary		Import Excel customer.
//	@Description	Import Excel customer.
//	@Produce		application/json
//	@Accept			multipart/form-data
//	@Tags			customers
//	@Param			file	formData	file	true	"Import Excel customer"
//	@Success		200		{object}	entity.JsonSuccess{data=string}"Data"
//	@Failure		400		{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404		{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500		{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/customers/import [post]
//	@Security		Bearer
func (controller *CustomerController) Import(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	file, err := ctx.FormFile("file")
	if err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	fmt.Println("cek", file)
	err = controller.customerService.Import(c, file)
	helper.ErrorPanic(err)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "Import Successful",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

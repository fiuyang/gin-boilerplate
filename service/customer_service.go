package service

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/tealeg/xlsx"
	"math"
	"mime/multipart"
	"scylla/entity"
	"scylla/model"
	"scylla/pkg/exception"
	"scylla/pkg/helper"
	"scylla/repository"
	"sync"
	"time"
)

type CustomerService interface {
	Create(ctx context.Context, request entity.CreateCustomerRequest) error
	CreateBatch(ctx context.Context, request entity.CreateCustomerBatchRequest)
	Update(ctx context.Context, request entity.UpdateCustomerRequest)
	DeleteBatch(ctx context.Context, request entity.DeleteBatchCustomerRequest)
	FindById(ctx context.Context, request entity.CustomerParams) (response entity.CustomerResponse)
	FindAll(ctx context.Context, dataFilter entity.CustomerQueryFilter) (response []entity.CustomerResponse)
	FindAllPaging(ctx context.Context, dataFilter entity.CustomerQueryFilter) (response []entity.CustomerResponse, paging entity.Meta)
	Export(ctx context.Context, dataFilter entity.CustomerQueryFilter) (string, error)
	Import(ctx context.Context, file *multipart.FileHeader) error
}

type CustomerServiceImpl struct {
	customerRepo repository.CustomerRepo
	validate     *validator.Validate
}

func NewCustomerServiceImpl(customerRepo repository.CustomerRepo, validate *validator.Validate) CustomerService {
	return &CustomerServiceImpl{
		customerRepo: customerRepo,
		validate:     validate,
	}
}

func (service *CustomerServiceImpl) Create(ctx context.Context, request entity.CreateCustomerRequest) error {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.Customer{
		Username: request.Username,
		Email:    request.Email,
		Phone:    request.Phone,
		Address:  request.Address,
	}

	err = service.customerRepo.Insert(ctx, dataset)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	return nil
}

func (service *CustomerServiceImpl) CreateBatch(ctx context.Context, request entity.CreateCustomerBatchRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	var customers []model.Customer
	for _, req := range request.Customers {
		customer := model.Customer{
			Username: req.Username,
			Email:    req.Email,
			Phone:    req.Phone,
			Address:  req.Address,
		}
		customers = append(customers, customer)
	}

	batchSize := len(request.Customers)

	err = service.customerRepo.InsertBatch(ctx, customers, batchSize)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}
}

func (service *CustomerServiceImpl) Update(ctx context.Context, request entity.UpdateCustomerRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset, err := service.customerRepo.FindById(ctx, request.ID)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	dataset.Username = request.Username
	dataset.Email = request.Email
	dataset.Phone = request.Phone
	dataset.Address = request.Address

	service.customerRepo.Update(ctx, dataset)
}

func (service *CustomerServiceImpl) DeleteBatch(ctx context.Context, request entity.DeleteBatchCustomerRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	service.customerRepo.DeleteBatch(ctx, request.ID)
}

func (service *CustomerServiceImpl) FindById(ctx context.Context, request entity.CustomerParams) (response entity.CustomerResponse) {
	result, err := service.customerRepo.FindById(ctx, request.CustomerId)

	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	helper.Automapper(result, &response)
	return response
}

func (service *CustomerServiceImpl) FindAll(ctx context.Context, dataFilter entity.CustomerQueryFilter) (response []entity.CustomerResponse) {
	result, err := service.customerRepo.FindAll(ctx, dataFilter)

	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	for _, row := range result {
		var res entity.CustomerResponse
		helper.Automapper(row, &res)
		response = append(response, res)
	}
	return response
}

func (service *CustomerServiceImpl) FindAllPaging(ctx context.Context, dataFilter entity.CustomerQueryFilter) (response []entity.CustomerResponse, paging entity.Meta) {

	result := service.customerRepo.FindAllPaging(ctx, dataFilter)

	for _, value := range result {
		var res entity.CustomerResponse
		helper.Automapper(value, &res)

		response = append(response, res)
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	if dataFilter.Page == 0 {
		dataFilter.Page = 1
	}

	var total int
	if len(result) > 0 {
		total = len(result)
	} else {
		total = 0
	}

	paging.Page = dataFilter.Page
	paging.Limit = dataFilter.Limit
	paging.TotalData = total
	paging.TotalPage = int(math.Ceil(float64(total) / float64(dataFilter.Limit)))

	return response, paging
}

func (service *CustomerServiceImpl) Export(ctx context.Context, dataFilter entity.CustomerQueryFilter) (string, error) {
	// Create a new Excel file
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Customer")
	helper.ErrorPanic(err)

	headers := []string{"ID", "Username", "Email", "Phone", "Address", "CreatedAt"}
	headerRow := sheet.AddRow()

	// Define the style for the header row
	style := xlsx.NewStyle()
	style.Font = *xlsx.NewFont(12, "Calibri")
	style.Fill = *xlsx.NewFill("solid", "00FFFF00", "00FFFF00") // Yellow background

	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
		cell.SetStyle(style)
	}

	result, err := service.customerRepo.FindAll(ctx, dataFilter)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	dataStyle := xlsx.NewStyle()
	dataStyle.Font = *xlsx.NewFont(12, "Calibri")

	for _, value := range result {
		dataRow := sheet.AddRow()
		cell := dataRow.AddCell()
		cell.SetInt(int(value.ID))
		cell.SetStyle(dataStyle)

		cell = dataRow.AddCell()
		cell.Value = value.Username
		cell.SetStyle(dataStyle)

		cell = dataRow.AddCell()
		cell.Value = value.Email
		cell.SetStyle(dataStyle)

		cell = dataRow.AddCell()
		cell.Value = value.Phone
		cell.SetStyle(dataStyle)

		cell = dataRow.AddCell()
		cell.Value = value.Address
		cell.SetStyle(dataStyle)

		cell = dataRow.AddCell()
		cell.Value = value.CreatedAt
		cell.SetStyle(dataStyle)
	}

	// Save the Excel file
	timestamp := time.Now().Format("2006-01-02_150405")
	filePath := fmt.Sprintf("customer_%s.xlsx", timestamp)

	err = file.Save(filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func (service *CustomerServiceImpl) Import(ctx context.Context, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return exception.NewInternalServerError(err.Error())
	}
	defer src.Close()

	xlFile, err := xlsx.OpenReaderAt(src, file.Size)
	if err != nil {
		return exception.NewInternalServerError(err.Error())
	}

	sheet := xlFile.Sheets[0]

	// Create channels for error handling and synchronization
	errorChan := make(chan error)
	validationExcel := exception.ValidationExcel{}
	wg := sync.WaitGroup{}

	// Track unique fields to ensure no duplicates within the same file
	uniqueTracker := make(map[string]map[string]bool)

	for colIndex := range helper.UniqueCustomer {
		uniqueTracker[helper.UniqueCustomer[colIndex]] = make(map[string]bool)
	}

	for rowIndex, row := range sheet.Rows {
		if rowIndex == 0 {
			continue
		}

		wg.Add(1)

		go func(rowIndex int, row *xlsx.Row) {
			defer wg.Done()

			// Validate each cell dynamically
			for colIndex, errorMsg := range helper.ColumnCustomer {
				if colIndex < len(row.Cells) {
					cell := row.Cells[colIndex]
					if cell.String() == "" {
						fieldName := sheet.Rows[0].Cells[colIndex].String()
						validationExcel.Add(fieldName, rowIndex, errorMsg)
					}
				}
			}

			// If there are validation errors, skip further processing for this row
			if len(validationExcel.Errors) > 0 {
				return
			}

			// Check uniqueness within the file
			for colIndex, fieldName := range helper.UniqueCustomer {
				if colIndex < len(row.Cells) {
					cell := row.Cells[colIndex]
					if uniqueTracker[fieldName][cell.String()] {
						validationExcel.Add(fieldName, rowIndex, fmt.Sprintf("%s '%s' is not unique", fieldName, cell.String()))
						return
					}
					uniqueTracker[fieldName][cell.String()] = true
				}
			}

			// Check uniqueness in the database
			for colIndex, fieldName := range helper.UniqueCustomer {
				if colIndex < len(row.Cells) {
					cell := row.Cells[1]
					data := service.customerRepo.CheckColumnExists(ctx, fieldName, cell.String())
					if data != false {
						validationExcel.Add(fieldName, rowIndex+1, fmt.Sprintf("%s '%s' already taken", fieldName, cell.String()))
						return
					}
				}
			}

			var customer model.Customer

			for i, cell := range row.Cells {
				switch i {
				case 0:
					customer.Username = cell.String()
				case 1:
					customer.Email = cell.String()
				case 2:
					customer.Phone = cell.String()
				case 3:
					customer.Address = cell.String()
				}
			}

			service.customerRepo.Insert(ctx, customer)

		}(rowIndex, row)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	for err := range errorChan {
		return exception.NewInternalServerError(err.Error())
	}

	if len(validationExcel.Errors) > 0 {
		return &validationExcel
	}

	return nil
}

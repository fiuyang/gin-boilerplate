package service

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/tealeg/xlsx"
	"mime/multipart"
	"scylla/dto"
	"scylla/entity"
	"scylla/pkg/exception"
	"scylla/pkg/helper"
	"scylla/pkg/utils"
	"scylla/repository"
	"strings"
	"sync"
	"time"
)

type UserService interface {
	Create(ctx context.Context, request dto.CreateUserRequest)
	Update(ctx context.Context, request dto.UpdateUserRequest)
	DeleteBatch(ctx context.Context, request dto.DeleteBatchUserRequest)
	FindAll(ctx context.Context, dataFilter dto.UserQueryFilter) (response []dto.UserResponse)
	FindById(ctx context.Context, params dto.UserParams) (response dto.UserResponse)
	Export(ctx context.Context, dataFilter dto.UserQueryFilter) (string, error)
	Import(ctx context.Context, file *multipart.FileHeader) error
}

type UserServiceImpl struct {
	userRepo repository.UserRepo
	validate *validator.Validate
}

func NewUserServiceImpl(userRepo repository.UserRepo, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
		validate: validate,
	}
}

func (service *UserServiceImpl) Create(ctx context.Context, request dto.CreateUserRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	hashedPassword, err := utils.HashPassword(request.Password)
	helper.ErrorPanic(err)

	dataset := entity.User{
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
	}

	err = service.userRepo.Insert(ctx, dataset)
	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}

}

func (service *UserServiceImpl) Update(ctx context.Context, request dto.UpdateUserRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset, err := service.userRepo.FindById(ctx, request.ID)
	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}

	if dataset.Password != "" {
		hashedPassword, err := utils.HashPassword(dataset.Password)
		helper.ErrorPanic(err)
		dataset.Password = hashedPassword
	}
	dataset.Username = request.Username
	dataset.Email = request.Email

	err = service.userRepo.Update(ctx, dataset)
	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}
}

func (service *UserServiceImpl) DeleteBatch(ctx context.Context, request dto.DeleteBatchUserRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	err = service.userRepo.DeleteBatch(ctx, request.ID)
	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}
}

func (service *UserServiceImpl) FindAll(ctx context.Context, dataFilter dto.UserQueryFilter) (response []dto.UserResponse) {
	result, err := service.userRepo.FindAll(ctx, dataFilter)

	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}

	for _, row := range result {
		var res dto.UserResponse
		helper.Automapper(row, &res)
		response = append(response, res)
	}
	return response
}

func (service *UserServiceImpl) FindById(ctx context.Context, params dto.UserParams) (response dto.UserResponse) {
	result, err := service.userRepo.FindById(ctx, params.UserId)

	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}

	helper.Automapper(result, &response)
	return response
}

func (service *UserServiceImpl) Export(ctx context.Context, dataFilter dto.UserQueryFilter) (string, error) {
	// Create a new Excel file
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Users")
	helper.ErrorPanic(err)

	headers := []string{"ID", "Username", "Email", "CreatedAt", "UpdatedAt"}
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

	users, err := service.userRepo.FindAll(ctx, dataFilter)
	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}

	dataStyle := xlsx.NewStyle()
	dataStyle.Font = *xlsx.NewFont(12, "Calibri")

	for _, user := range users {
		dataRow := sheet.AddRow()
		cell := dataRow.AddCell()
		cell.SetInt(int(user.ID))
		cell.SetStyle(dataStyle)

		cell = dataRow.AddCell()
		cell.Value = user.Username
		cell.SetStyle(dataStyle)

		cell = dataRow.AddCell()
		cell.Value = user.Email
		cell.SetStyle(dataStyle)

		cell = dataRow.AddCell()
		cell.Value = user.CreatedAt.Format("2006-01-02")
		cell.SetStyle(dataStyle)

		cell = dataRow.AddCell()
		cell.Value = user.UpdatedAt.Format("2006-01-02")
		cell.SetStyle(dataStyle)
	}

	// Save the Excel file
	timestamp := time.Now().Format("2006-01-02_150405")
	filePath := fmt.Sprintf("user_%s.xlsx", timestamp)

	err = file.Save(filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func (service *UserServiceImpl) Import(ctx context.Context, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}
	defer src.Close()

	xlFile, err := xlsx.OpenReaderAt(src, file.Size)
	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}

	sheet := xlFile.Sheets[0]

	// Create channels for error handling and synchronization
	errorChan := make(chan error)
	excelValidation := exception.NewExcelValidationError{}
	wg := sync.WaitGroup{}

	uniqueTracker := make(map[string]map[string]bool)

	for _, rule := range helper.RulesExcelUser {
		fieldName := strings.Split(rule, ",")[0]
		uniqueTracker[fieldName] = make(map[string]bool)
	}

	for rowIndex, row := range sheet.Rows {
		if rowIndex == 0 {
			continue
		}

		wg.Add(1)

		go func(rowIndex int, row *xlsx.Row) {
			defer wg.Done()

			rowErrors := map[string][]string{}

			// Validate each cell dynamically
			for colIndex, rule := range helper.RulesExcelUser {
				fields := strings.Split(rule, ",")
				fieldName := fields[0]
				cell := row.Cells[colIndex]
				rules := fields[1:]

				for _, r := range rules {
					switch r {
					case "required":
						if cell.String() == "" {
							rowErrors[fieldName] = append(rowErrors[fieldName], fmt.Sprintf("%s row %d is required", fieldName, rowIndex+1))
						}
					case "unique":
						if uniqueTracker[fieldName][cell.String()] {
							rowErrors[fieldName] = append(rowErrors[fieldName], fmt.Sprintf("%s '%s' is not unique row %d", fieldName, cell.String(), rowIndex+1))
							return
						}
						uniqueTracker[fieldName][cell.String()] = true
					}
				}
			}

			// If there are validation errors, skip further processing for this row
			if len(rowErrors) > 0 {
				for field, errs := range rowErrors {
					for _, err := range errs {
						excelValidation.AddHandler(field, rowIndex+1, err)
					}
				}
				return
			}

			// Check unique in the database
			for colIndex, rule := range helper.RulesExcelUser {
				fields := strings.Split(rule, ",")
				fieldName := fields[0]
				rules := fields[1:]
				for _, r := range rules {
					if r == "unique" {
						cell := row.Cells[colIndex]
						data := service.userRepo.CheckColumnExists(ctx, fieldName, cell.String())
						if data != false {
							excelValidation.AddHandler(fieldName, rowIndex+1, fmt.Sprintf("%s '%s' already taken", fieldName, cell.String()))
							return
						}
					}
				}
			}

			var user entity.User

			for i, cell := range row.Cells {
				switch i {
				case 0:
					user.Username = cell.String()
				case 1:
					user.Email = cell.String()
				case 2:
					hashedPassword, err := utils.HashPassword(cell.String())
					if err != nil {
						errorChan <- fmt.Errorf("error hashing password in row %d: %v", rowIndex, err)
						return
					}
					user.Password = hashedPassword
				}
			}

			service.userRepo.Insert(ctx, user)

		}(rowIndex, row)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	for err := range errorChan {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}

	if len(excelValidation.Errors) > 0 {
		return &excelValidation
	}

	return nil
}

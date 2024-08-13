package service

import (
	"context"
	"errors"
	"fmt"
	"scylla/dto"
	"scylla/entity"
	"scylla/pkg/config"
	"scylla/pkg/exception"
	"scylla/pkg/helper"
	"scylla/pkg/utils"
	"scylla/repository"
	"time"

	"github.com/go-playground/validator/v10"
)

type AuthService interface {
	Login(ctx context.Context, request dto.LoginRequest) (token string, err error)
	Register(ctx context.Context, request dto.CreateUserRequest)
	Logout(token string) error
	ForgotPassword(ctx context.Context, request dto.ForgotPasswordRequest) (string, error)
	CheckOtp(ctx context.Context, request dto.CheckOtpRequest) (string, error)
	ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) (string, error)
}

type AuthServiceImpl struct {
	config        *config.Config
	userRepo      repository.UserRepo
	passResetRepo repository.PassResetRepo
	validate      *validator.Validate
}

func NewAuthServiceImpl(config *config.Config, userRepo repository.UserRepo, passResetRepo repository.PassResetRepo, validate *validator.Validate) AuthService {
	return &AuthServiceImpl{
		config:        config,
		userRepo:      userRepo,
		passResetRepo: passResetRepo,
		validate:      validate,
	}
}

func (service *AuthServiceImpl) Login(ctx context.Context, request dto.LoginRequest) (token string, err error) {
	error := service.validate.Struct(request)
	helper.ErrorPanic(error)

	data, err := service.userRepo.FindByColumns(ctx, []string{"email"}, []any{request.Email})
	if err != nil {
		return "", exception.NewBadRequestHandler("email or password is wrong")
	}

	err = utils.VerifyPassword(data.Password, request.Password)
	if err != nil {
		return "", exception.NewBadRequestHandler("email or password is wrong")
	}

	// Generate Token
	token, err = utils.GenerateToken(service.config.Jwt.Expire, data, service.config.Jwt.Secret)
	helper.ErrorPanic(err)

	return token, nil
}

func (service *AuthServiceImpl) Register(ctx context.Context, request dto.CreateUserRequest) {
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

func (service *AuthServiceImpl) Logout(token string) error {
	if token == "" {
		return errors.New("empty token")
	}

	err := utils.AddToBlacklist(token)
	helper.ErrorPanic(err)
	return nil
}

func (service *AuthServiceImpl) ForgotPassword(ctx context.Context, request dto.ForgotPasswordRequest) (string, error) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	data, err := service.userRepo.FindByColumns(ctx, []string{"email"}, []any{request.Email})
	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}

	otp := utils.GenerateOTP(4)
	if err != nil {
		return "", errors.New("failed to generate token otp")
	}

	dataset := entity.PasswordReset{
		Email:     data.Email,
		Otp:       otp,
		CreatedAt: time.Now().Add(time.Minute * 5),
	}

	err = service.passResetRepo.Insert(ctx, dataset)
	if err != nil {
		return "", exception.NewInternalServerErrorHandler(err.Error())
	}

	emailData := utils.EmailData{
		Otp:     otp,
		Email:   data.Email,
		Subject: "Reset Password",
	}

	utils.SendEmail(&data, &emailData, "resetPassword.html")

	return fmt.Sprintf("%d", otp), nil
}

func (service *AuthServiceImpl) CheckOtp(ctx context.Context, request dto.CheckOtpRequest) (string, error) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	data, err := service.passResetRepo.FindByColumns(ctx, []string{"otp"}, []any{request.Otp})
	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}

	if request.Otp != data.Otp {
		return "", errors.New("invalid otp")
	}

	if time.Now().After(data.CreatedAt) {
		return "", errors.New("otp has expired")
	}

	return "Otp Valid", nil
}

func (service *AuthServiceImpl) ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) (string, error) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	data, err := service.passResetRepo.FindByColumns(ctx, []string{"otp"}, []any{request.Otp})
	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}

	if request.Otp != data.Otp {
		return "", errors.New("invalid otp")
	}

	if time.Now().After(data.CreatedAt) {
		return "", errors.New("otp has expired")
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	helper.ErrorPanic(err)

	dataset := entity.User{
		Email:    data.Email,
		Password: hashedPassword,
	}

	err = service.userRepo.Update(ctx, dataset)
	if err != nil {
		return "", exception.NewNotFoundHandler(err.Error())
	}

	err = service.passResetRepo.DeleteByColumns(ctx, []string{"otp"}, []any{data.Otp})
	if err != nil {
		return "", exception.NewNotFoundHandler(err.Error())
	}

	return "Reset Password Successful", nil
}

package controllers

import (
	"context"
	"database/sql"
	"os"
	v1 "protoUserService/pkg/api/v1"
	"protoUserService/pkg/services/api/v1/models"
	"protoUserService/pkg/services/api/v1/validation"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// UserServices implemented on version 1 proto interface
type userServices struct {
	db *sql.DB
}

// New sending otp services create sending otp service
func NewUserServicesService(db *sql.DB) v1.UserServicesServer {
	return &userServices{db: db}
}

// Checking Api version
func (s *userServices) CheckApi(api string) error {
	if len(api) > 0 {
		if os.Getenv("API_VERSION") != api {
			return status.Errorf(codes.Unimplemented, "Unsupported API Version: Service API implement using '%s', but asked for '%s'", os.Getenv("API_VERSION"), api)
		}
	}
	return nil
}

// Function login user account
func (s *userServices) LoginAccount(ctx context.Context, req *v1.LoginAccountRequest) (*v1.LoginAccountResponse, error) {
	// Get data IP Address
	// dataIp, _ := peer.FromContext(ctx)
	// timeRequest := time.Now().Format("2006-01-02 15:04:05")
	// var code, status, message, token, fullName string
	// var isActive bool

	// message, err := validation.LoginRequest(dataIp.Addr.String(), timeRequest, req.Api, req.NumberPhone, req.Password)

	// if err == false {
	// 	message = message
	// 	status = "Field is not match requirement"
	// 	code = "422"
	// } else {
	// 	if err := s.CheckApi(req.Api); err != nil {
	// 		return nil, err
	// 	} else {
	// code, status, message, token, fullName, isActive = models.Login
	// 	}
	// }

	return &v1.LoginAccountResponse{
		Code:    "00",
		Status:  "berhasil",
		Message: "Berhasil",
		Data: &v1.DataResponseAccount{
			Token:      "dads",
			IsActive:   true,
			FullName:   "Dasd",
			LoggedTime: "asdasd",
		},
	}, nil
}

// Func register account
func (s *userServices) RegisterAccount(ctx context.Context, req *v1.RegisterAccountRequest) (*v1.RegisterAccountResponse, error) {
	// Get data IP Address
	dataIp, _ := peer.FromContext(ctx)
	timeRequest := time.Now().Format("2006-01-02 15:04:05")
	var code, status, message, token, fullName string
	var isActive bool

	message, err := validation.RegistrationRequest(req.Api, req.NumberPhone, req.FullName, req.Password, timeRequest, dataIp.Addr.String())

	if err == false {
		message = message
		status = "Field is not match requirement"
		code = "422"
	} else {
		if err := s.CheckApi(req.Api); err != nil {
			return nil, err
		} else {
			code, status, message, token, fullName, isActive = models.RegisterAccount(dataIp.Addr.String(), req.NumberPhone, req.FullName, req.Password, timeRequest, s.db, ctx)
		}
	}

	return &v1.RegisterAccountResponse{
		Code:    code,
		Status:  status,
		Message: message,
		Data: &v1.DataResponseAccount{
			Token:      token,
			FullName:   fullName,
			IsActive:   isActive,
			LoggedTime: timeRequest,
		},
	}, nil
}

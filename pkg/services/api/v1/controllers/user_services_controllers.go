package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	v1 "protoUserService/pkg/api/v1"
	"protoUserService/pkg/services/api/v1/global"
	"protoUserService/pkg/services/api/v1/models"
	"protoUserService/pkg/services/api/v1/validation"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

// Function to check token valid or not
func Auth(data string) (string, bool) {
	tokenize := &global.Tokenization{}

	tkn, err := jwt.ParseWithClaims(data, tokenize, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN")), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "Invalid signature", false
		}
		return "Bas Request", false
	}

	if !tkn.Valid {
		return "Invalid Token", false
	} else {
		fmt.Println(tokenize)
		return "Token valid", true
	}
}

// Func percobaan header
func (s *userServices) DataCoba(ctx context.Context, req *v1.DataRequest) (*v1.DataResponse, error) {
	headers, _ := metadata.FromIncomingContext(ctx)
	data := headers["authorization"]
	fmt.Println(data)

	decodedData, _ := models.Decrypt(data[0])
	fmt.Println(decodedData)

	token, status := Auth(decodedData)

	if status == false {
		return &v1.DataResponse{
			Output: token,
		}, nil
	} else {
		return &v1.DataResponse{
			Output: token,
		}, nil
	}
}

// Func register account
func (s *userServices) RegisterAccount(ctx context.Context, req *v1.RegisterAccountRequest) (*v1.RegisterAccountResponse, error) {
	// Get data IP Address
	dataIp, _ := peer.FromContext(ctx)
	timeRequest := time.Now().Format("2006-01-02 15:04:05")
	var code, status, message, token, fullName string
	var isActive bool

	message, err := validation.RegistrationRequest(req.Api, req.NumberPhone, req.Username, req.FullName, req.Password, timeRequest, dataIp.Addr.String())

	if err == false {
		message = message
		status = "Field is not match requirement"
		code = "422"
	} else {
		if err := s.CheckApi(req.Api); err != nil {
			return nil, err
		} else {
			code, status, message, token, fullName, isActive = models.RegisterAccount(dataIp.Addr.String(), req.NumberPhone, req.Username, req.FullName, req.Password, timeRequest, s.db, ctx)
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

package models

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"protoUserService/pkg/services/api/v1/global"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mergermarket/go-pkcs7"
	"github.com/minio/highwayhash"
)

// Set global environment variablea
var conf *global.Configuration
var messageError map[int]global.MessageError
var level, cases, fatal string

// Function initialization
func init() {
	conf = global.New()
	messageError = global.GetMessageError()
}

// function to encrypt using AES 256
func Encrypt(data string) (string, error) {
	key := []byte(conf.KeyAes)
	plainText := []byte(data)
	plainText, err := pkcs7.Pad(plainText, aes.BlockSize)

	if err != nil {
		err := fmt.Errorf(`text: "%s" block size salah`, plainText)
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(crand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return fmt.Sprintf("%x", cipherText), nil
}

// Function to decrypt
func Decrypt(encrypted string) (string, error) {
	key := []byte(conf.KeyAes)
	cipherText, _ := hex.DecodeString(encrypted)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(cipherText) < aes.BlockSize {
		panic("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		panic("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, _ = pkcs7.Unpad(cipherText, aes.BlockSize)
	return fmt.Sprintf("%s", cipherText), nil
}

// Function Generate random number
func getRandomString() string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	const charset = "1234567890" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

// Function to create a JWT Token
func JwtTokenCreate(token *global.Tokenization) (string, error) {
	tokenAuth := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token)
	tokenNew, _ := tokenAuth.SignedString([]byte(conf.Token))

	return tokenNew, nil
}

// Function password hashing
func hashPasswordData(password string) string {
	passwordHash := highwayhash.Sum([]byte(password), []byte(conf.KeyPass))
	pass := hex.EncodeToString(passwordHash[:])

	return fmt.Sprintf(pass)
}

// Function for create Token JWT
func createToken(userId, timeRequest, ipAddress, numberPhone string) string {
	tokenJwt := &global.Tokenization{
		UserId:      userId,
		Time:        timeRequest,
		Ip:          ipAddress,
		NumberPhone: numberPhone,
	}

	tokenAuth, _ := JwtTokenCreate(tokenJwt)
	encrypted, _ := Encrypt(tokenAuth)

	return encrypted
}

// Function for register account
func RegisterAccount(ipAddress, numberPhone, fullname, password, timeRequest string, connection *sql.DB, ctx context.Context) (code, status, message, token, fullName string, isActive bool) {
	var count int

	// check username and numberPhone
	checkUsername := global.GenerateQueryForUser(map[string]string{
		"numberPhone": numberPhone,
		"email":       "",
	})

	fmt.Println(checkUsername)
	rows := connection.QueryRowContext(ctx, checkUsername)

	err := rows.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}

	if count == 0 {
		pass := hashPasswordData(password)
		userId := getRandomString()

		token = createToken(userId, timeRequest, ipAddress, numberPhone)

		go func() {
			sql := `INSERT INTO "user" ("id", "phoneNumber", "fullName", "password", "status", "createdAt") VALUES ($1, $2, $3, $4, $5, $6)`

			_, err := connection.Query(sql, userId, numberPhone, fullname, pass, 1, timeRequest)

			if err != nil {
				fmt.Println(err)
			}
		}()

		code = messageError[00].Code
		message = messageError[00].Message
		status = "Registration Success"
		isActive = true
	} else {
		code = messageError[422].Code
		message = messageError[422].Message
		status = "Failed registration, numberPhone has been taken"
		token = ""
		fullName = fullname
		isActive = false
	}

	return code, status, message, token, fullname, isActive
}

// Login request

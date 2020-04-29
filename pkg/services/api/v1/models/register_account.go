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

// Set global environment variable
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
		err := fmt.Errorf(`plainText: "%s" has the wrong block size`, plainText)
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

// Function for register account
func RegisterAccount(ipAddress, username, password, timeRequest string, connection *sql.DB, ctx context.Context) (code, status, message, dataUsername string) {
	var count int

	// check username
	checkUsername := global.GenerateQueryForUser(map[string]string{
		"username": username,
	})

	rows := connection.QueryRowContext(ctx, checkUsername)

	err := rows.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}

	if count == 0 {
		passwordHash := highwayhash.Sum([]byte(password), []byte(conf.KeyPass))
		pass := hex.EncodeToString(passwordHash[:])
		userId := getRandomString()

		go func() {
			sql := `INSERT INTO "users" ("userId", "username", "password", "status") VALUES ($1, $2, $3, $4)`

			_, err := connection.Query(sql, userId, username, pass, "isActive")

			if err != nil {
				fmt.Println(err)
			}
		}()

		code = messageError[00].Code
		message = messageError[00].Message
		dataUsername = username
		status = "Registration Success"
	} else {
		code = messageError[422].Code
		message = "Failed registration, username has been taken"
		status = messageError[422].Message
		dataUsername = username
	}

	return code, status, message, dataUsername
}

// Func login
func LoginAccount(ipAddress, username, password, timeRequest string, connection *sql.DB, ctx context.Context) (code, status, message, dataUsername, token string) {
	passwordHash := highwayhash.Sum([]byte(password), []byte(conf.KeyPass))
	pass := hex.EncodeToString(passwordHash[:])

	var userId string

	// check username
	checkUsername := global.GenerateQueryForLogin(map[string]string{
		"username": username,
		"password": pass,
	})

	rows := connection.QueryRowContext(ctx, checkUsername)

	err := rows.Scan(&userId)
	if err != nil {
		fmt.Println(err)
	}

	if userId != "" {
		tokenJwt := &global.Tokenization{
			UserId:   userId,
			Password: pass,
			Time:     timeRequest,
			Ip:       ipAddress,
			Username: username,
		}

		tokenAuth, _ := JwtTokenCreate(tokenJwt)
		encrypted, _ := Encrypt(tokenAuth)

		token = encrypted
		code = "00"
		status = "Success"
		message = "Login success"
		dataUsername = username

	} else {
		token = "nil"
		code = "05"
		status = "Failed Login"
		message = "Invalid username or password"
		dataUsername = username
	}

	return code, status, message, dataUsername, token
}

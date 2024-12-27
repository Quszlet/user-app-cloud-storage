package service

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Quszlet/authentication_service/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCtNz9VKF1o/F0+
glrhBM7EXFHI8fMCl30YOGKr1byR/XWXui7svJ0u0ljmXKou2WytiPS8gDUkrb2d
Kbk00gHtrcsTArNQLxUblDqn7/pmynkz6Xn+fbGHdVCCRwXqOxdHLqjY3CpHb/ZT
/ouBR8rW64DEv0+3h4jJt9Y6CHFD1Qs+Vs2N5cCfoapVSPyW0PEKjXylfPTo4X+a
8zyuY7pNSvlW8LBI7/CbA41DSfx763crMGr8sGhF4ZgJoMiVtaFrCbmXdWQpkxzA
LGofrJUlGJj1EeD/FvGT7IPnpJz5lhQQH4aEOto6mkkdSEihNG8jultFfVzHH2aD
impN5cN5AgMBAAECggEAGUsyAltsmIIRk6kxYB51mxLoFnHOTJiWeczyC6mgaakb
XNahN4yrX0T0Gl95snGpfwW6xzPpjruYLrUDCIWKJoz0VIhWauUVLfvNPOy8Ifp7
DGuRlun/3InzAbMhV/zc/1X+7fvoaJoe3FEVSgGTyyKEoTZWi7RM8zfXHo5RC9h4
jyHppBmRINBfQDNUdwV31ooevwBqJxkYaWefVOM7tys7d2aHFbaKxvMmq2lTD8pQ
D666xcmt2CZXAemEIWHlMqKeFVXXnURVMdcoeZmbfJfwz+1ffG6AtJZqj4GY+JfU
sw00d3JBYKUmICjKQL4awMWoMru+YHiOF0DiK66SaQKBgQDngDQ2xJHU/g7dfRvs
ORYgHs7uTBhfoxPt8rtFgNBwP+iDgdNHQmK6FcEbw9cdRXJ8ERnrPeyaqXWoM54i
B13lNMcG6Si5bC44S0OhC3q14K9IXNVyQ2HjvDR4eui6y+X6EpFy+I6AATv/KJWm
je4Z1gl2DPhdlm6V2/vqIiwQNwKBgQC/i/4UufDbrYR11ldrJam5+Az0E8H20TEl
hOhIl5kofyelblcths1mhivE8ZDEHJhWDYwwH1/tcssVVUYqyQg3yLJH5SO2sbt/
VgezSK03HBZZ9fbKYoWs+kAQ1JL/EEtutjAWvPs+YTZNqnfTM4angqVCg+c/YpTt
HzYUZ9cRzwKBgEfaFdqt1imKlSiPrTv7V++uQEHcInCEmCnxfciLu6YrX8p5YA8s
/qGNIPuyJDE0ndz+HdJSzP4P/LGxG7KqIK2EXQW7FmW+uvXD/oCcpICQ9TZ7gdBO
M7LQdmSymUto/79HRheuJ+R4/ZsriI9CXBVuxk76pZe+miIvPhgkdRKvAoGASnKe
zm92retDEIm+cGazERTX9AW53bRw5aRCo/RIEvVY83NvbsJ2EuMTH/jDy3VRwCCD
3DBVmHSFekUqgHaiOwxGPqtQtjFeLp/BXm1g5YqFJXHz+bVRP2oEfIYinAA5UU5+
YlgRTq157bXT3MKqE3EzyiZ6OqiWOZNn8YZqkQECgYAwRpuFYtaGpHYUhsS6zp8j
uIrAgQMXFdHlZfRBujusix/I8H7pTHTswiMU1l9XlqoNAKmSbuI+rNEFlkZZqRdp
DeWsl+4pMSjJuXTgKrOZLw9X+XwkD+wsDJCOJ32lDahC/jzPeTo/mlK24tACI7f0
LM/OkyFpOOjTg8ZY1h+JMw==
-----END RSA PRIVATE KEY-----`)

type AuthService struct {
	TimeOut int64
}

func (as *AuthService) SignIn(Username, Password string) (string, error) {
	user, err := as.GetUser(Username, Password)

	if err != nil {
		return "", fmt.Errorf("error get user: %v", err)
	}

	if user.Login == Username && user.Password == Password {
		token, err := as.GenerateJWTToken(user.Id, user.Login)
		if err != nil {
			return "", fmt.Errorf("error generate JWT token: %v", err)
		}

		return token, nil
	}

	return "", fmt.Errorf("login or password is not correct")
}

func (as *AuthService) GetUser(Username, Password string) (*models.User, error) {
	url := fmt.Sprintf("http://user_service:8010/users/%s", Username) // Микросервис пользователей

	resp, err := http.Get(url)
	if err != nil {
		return &models.User{}, err
	}

	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &models.User{}, fmt.Errorf("error with read body. %v", err)
	}

	responseData := models.User{}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return &models.User{}, fmt.Errorf("error with deserialize body. %s", body)
	}

	return &responseData, nil
}

func (as *AuthService) GenerateJWTToken(UserId int, Username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		UserID:   UserId,
		Username: Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	block, _ := pem.Decode(privateKey)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", fmt.Errorf("failed to decode PEM block containing private key")
	}

	var err error
	privateKeyRSA, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse RSA private key: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKeyRSA)
	return tokenString, err
}

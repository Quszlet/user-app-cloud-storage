package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Quszlet/authentication_service/internal/models"
)

type AuthService struct {
	TimeOut time.Time
}

func (as *AuthService) SignIn(Username, Password string) (string, error) {
	username, err := as.GetUser(Username, Password)

	if err != nil {
		return "", err 
	}

	token, err := as.GenerateJWTToken(username)
	
}

func (as *AuthService) GetUser(Username, Password string) (string, error) {
	url := fmt.Sprintf("http://127.0.0.1:8080/users/%s", Username) // Микросервис пользователей

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error with read body: %v", err)
	}

	responseData := models.User{}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "", fmt.Errorf("Error with deserialize body: %v", err)
	}

}

func (as *AuthService) GenerateJWTToken(Username string) (string, error) {

}
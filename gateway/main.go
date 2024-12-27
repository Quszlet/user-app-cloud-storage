package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println(r.URL.Path)
		if r.URL.Path == "/users/create" {
			next.ServeHTTP(w, r)
			return
		}

		// Извлечение токена из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Валидация токена
		claims, err := ValidateToken(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// Логирование аутентифицированного пользователя
		log.Printf("Authenticated user: %s and id: %d", claims.Username, claims.UserID)

		// Пропуск запроса дальше
		next.ServeHTTP(w, r)
	})
}

func ValidateToken(tokenString string) (*Claims, error) {
	publicKey := []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArTc/VShdaPxdPoJa4QTO
xFxRyPHzApd9GDhiq9W8kf11l7ou7LydLtJY5lyqLtlsrYj0vIA1JK29nSm5NNIB
7a3LEwKzUC8VG5Q6p+/6Zsp5M+l5/n2xh3VQgkcF6jsXRy6o2NwqR2/2U/6LgUfK
1uuAxL9Pt4eIybfWOghxQ9ULPlbNjeXAn6GqVUj8ltDxCo18pXz06OF/mvM8rmO6
TUr5VvCwSO/wmwONQ0n8e+t3KzBq/LBoReGYCaDIlbWhawm5l3VkKZMcwCxqH6yV
JRiY9RHg/xbxk+yD56Sc+ZYUEB+GhDraOppJHUhIoTRvI7pbRX1cxx9mg4pqTeXD
eQIDAQAB
-----END PUBLIC KEY-----`)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(publicKey)
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// createReverseProxy создает прокси-сервер для перенаправления запросов на указанный URL
func createReverseProxy(target string) *httputil.ReverseProxy {
	url, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Ошибка парсинга URL: %v", err)
	}
	return httputil.NewSingleHostReverseProxy(url)
}

// handleProxy обрабатывает запросы, перенаправляя их на целевой сервис
func handleProxy(proxy *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Перенаправление запроса %s %s", r.Method, r.URL.Path)
		proxy.ServeHTTP(w, r)
	}
}

func main() {

	// TODO: вынести порты в окружение
	userService := "http://user_service:8010"
	fileSystemSevice := "http://file_system_service:8020"
	authenticationService := "http://authentication_service:8030"

	// Создаем прокси для каждого сервиса
	proxyUserService := createReverseProxy(userService)
	proxyFileSystemService := createReverseProxy(fileSystemSevice)
	proxyAuthenticationService := createReverseProxy(authenticationService)

	// Создаем маршрутизатор Mux
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Gateway работает. Доступные маршруты: /users/, /files/, /directories/")
	})

	checkValidation := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Токен валиден")
	})


	// Пример маршрутов с использованием прокси и middleware
	r.PathPrefix("/users").Handler(authMiddleware(handleProxy(proxyUserService)))
	r.PathPrefix("/files").Handler(authMiddleware(handleProxy(proxyFileSystemService)))
	r.PathPrefix("/directories").Handler(authMiddleware(handleProxy(proxyFileSystemService)))
	r.PathPrefix("/validation").Handler(authMiddleware(checkValidation))
	r.PathPrefix("/authentication").Handler(handleProxy(proxyAuthenticationService))


	// Корневой маршрут для проверки работоспособности
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Gateway работает. Доступные маршруты: /user/ и /file/")
	})


	corsMiddleware := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"}, // Укажите домены или используйте "*" для всех
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    })

	handler := corsMiddleware.Handler(r)

	log.Printf("Запуск Gateway на порту %s...", os.Getenv("PORT"))
	if err := http.ListenAndServe(os.Getenv("PORT"), handler); err != nil {
		log.Fatalf("Ошибка запуска Gateway: %v", err)
	}
}

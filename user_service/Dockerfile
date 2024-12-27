# Используем официальный базовый образ Go
FROM golang:1.23-alpine

RUN apk add --no-cache curl

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем все файлы проекта
COPY . .

# Собираем приложение из директории cmd/app
RUN go build -o user_service ./cmd/app

# Экспонируем порт, который использует приложение
EXPOSE 8010

# Команда для запуска приложения
CMD ["./user_service"]

FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Компилируем Go-приложение
RUN CGO_ENABLED=1 GOOS=linux go build -o /scheduler-app -a -ldflags '-linkmode external -extldflags "-static"' cmd/main.go && ls -la /scheduler-app

# Минимальный базовый образ для запуска приложения
FROM ubuntu:latest

# Устанавливаем необходимые пакеты
RUN apt-get update && apt-get install -y sqlite3

# Создаем рабочую директорию
WORKDIR /myapp

# Копируем скомпилированный бинарный файл из предыдущего этапа
COPY --from=builder /scheduler-app /myapp/scheduler-app

# Копируем директорию web
COPY ./web /myapp/web

# Указываем порт, который будет использоваться веб-сервером
EXPOSE 7545

# Определяем переменные окружения
COPY .env /myapp

# Команда для запуска веб-сервера

CMD [ "./scheduler-app" ]
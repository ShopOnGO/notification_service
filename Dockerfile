FROM golang:1.23.3 AS builder

WORKDIR /notification

# Устанавливаем pg_isready и очищаем кеш
RUN apt-get update && apt-get install -y postgresql-client \
    && rm -rf /var/lib/apt/lists/* && apt-get clean

# Отключаем CGO для статической компиляции
 ENV CGO_ENABLED=0

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download && go mod verify

# Копируем весь код
COPY . .

# Компилируем бинарник
RUN go build -o /notification/notification_service ./cmd/main.go



# Второй этап: финальный образ (без лишних инструментов)
FROM alpine:latest

WORKDIR /notification

# Устанавливаем postgresql-client и dos2unix
RUN apk add --no-cache postgresql-client dos2unix

# COPY .env /notification/.env

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /notification/notification_service /notification/notification_service

# Копируем wait-for-db.sh и делаем исполняемым
COPY --from=builder /notification/wait-for-db.sh /notification/wait-for-db.sh
RUN chmod +x /notification/wait-for-db.sh

# Преобразуем формат строки в скрипте wait-for-db.sh в Unix-формат
RUN dos2unix /notification/wait-for-db.sh

# Запуск приложения
CMD ["/notification/notification_service"]

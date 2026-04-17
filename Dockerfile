FROM golang:1.26.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем именно из cmd/apiserver (как у вас в структуре проекта)
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/mifare ./cmd/apiserver

# ====================== Финальная стадия ======================
FROM nginx:alpine

RUN apk --no-cache add ca-certificates

RUN mkdir -p /app /etc/nginx/certs

# Копируем бинарник Go-приложения
COPY --from=builder /bin/mifare /app/mifare

# Конфиги Nginx и TLS
COPY nginx/nginx.conf /etc/nginx/nginx.conf
COPY certs/cert.pem certs/key.pem /etc/nginx/certs/

# Конфиги приложения и миграции (ОЧЕНЬ ВАЖНО!)
COPY configs /app/configs
COPY migrations /app/migrations

# Устанавливаем рабочую директорию — теперь все относительные пути работают
WORKDIR /app

RUN chmod +x /app/mifare

EXPOSE 8888

# Запускаем и Go-приложение, и nginx одновременно
CMD ["sh", "-c", "./mifare & nginx -g 'daemon off;'"]
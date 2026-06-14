FROM golang:1.26.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/mifare ./cmd/apiserver

FROM node:20-alpine AS frontend-builder
WORKDIR /frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install
COPY frontend/ ./
RUN npm run build

FROM nginx:alpine
RUN apk --no-cache add ca-certificates

COPY --from=builder /bin/mifare /app/mifare
COPY --from=frontend-builder /frontend/dist /usr/share/nginx/html
COPY nginx/nginx.conf /etc/nginx/nginx.conf
COPY certs/cert.pem certs/key.pem /etc/nginx/certs/
COPY configs /app/configs
COPY migrations /app/migrations

WORKDIR /app
RUN chmod +x /app/mifare

EXPOSE 8888

# Запускаем Nginx, а Go-сервер запускается в фоне
CMD ["sh", "-c", "./mifare & nginx -g 'daemon off;'"]
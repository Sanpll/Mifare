# ================== STAGE 1: Build Frontend ==================
FROM node:20-alpine AS frontend-builder
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ================== STAGE 2: Build Backend ==================
FROM golang:1.26.2 AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/mifare ./cmd/apiserver

# ================== STAGE 3: Final Image ==================
FROM nginx:alpine
RUN apk --no-cache add ca-certificates

# Копируем собранный React
COPY --from=frontend-builder /frontend/dist /usr/share/nginx/html

# Копируем Go бинарник
COPY --from=backend-builder /bin/mifare /app/mifare

# Конфиги
COPY nginx/nginx.conf /etc/nginx/nginx.conf
COPY certs/cert.pem certs/key.pem /etc/nginx/certs/
COPY configs /app/configs
COPY migrations /app/migrations

WORKDIR /app
RUN chmod +x /app/mifare

EXPOSE 8888

# Запускаем Go + Nginx
CMD ["sh", "-c", "./mifare & nginx -g 'daemon off;'"]
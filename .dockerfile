FROM golang:1.21-alpine AS builder

# Установка зависимостей для сборки с AVX
RUN apk add --no-cache gcc musl-dev linux-headers

WORKDIR /app
COPY . .

# Сборка с флагами оптимизации
RUN go build -tags=avx -ldflags="-s -w" -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/main .

# Убедимся, что процессор поддерживает AVX
CMD ["/root/main"]
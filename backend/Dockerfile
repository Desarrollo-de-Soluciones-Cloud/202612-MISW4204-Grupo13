# ============================================
# ESTAPA 1: Builder
# ============================================
FROM golang:1.22-alpine AS builder
# Instalar dependencias para compilación
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
# Copiar go.mod y go.sum PRIMERO (beneficia caché de Docker)
COPY go.mod go.sum ./
# Descargar dependencias
RUN go mod download
# Copiar código fuente
COPY . .
# Compilar binario (CGO_ENABLED=0 para imagen final más pequeña)
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api
# ============================================
# ESTAPA 2: Producción
# ============================================
FROM alpine:3.19
WORKDIR /root/
# Instalar certificado CA para HTTPS
RUN apk --no-cache add ca-certificates
# Copiar binario compilado del builder
COPY --from=builder /app/server .
# Variables de entorno para configuración
ENV PORT=8080
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=postgres
# Exponer puerto
EXPOSE 8080
# Usuario no root por seguridad
RUN adduser -D -g '' appuser
USER appuser
CMD ["./server"]
# Build stage

FROM hub.docker.humo.lab/nexus-repository-golang-alpine3.23 AS builder
ARG CI_JOB_TOKEN

WORKDIR /app

# Сначала копируем только файлы зависимостей для кэширования слоя
COPY go.mod go.sum ./
RUN go env -w GOPRIVATE=gitlab.humo.tj
RUN go mod download

# Затем копируем весь остальной код
COPY . .
RUN go mod tidy

RUN go build -o converterApi cmd/main.go

# Run stage

FROM hub.docker.humo.lab/nexus-repository-alpine

WORKDIR /app

COPY --from=builder /app/converterApi .
COPY --from=builder /app/internal/config/config.json ./internal/config/config.json

CMD ["./converterApi"]
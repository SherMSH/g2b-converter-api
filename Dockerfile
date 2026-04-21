# Build stage

FROM hub.docker.humo.lab/nexus-repository-golang-alpine3.23 AS builder
ARG CI_JOB_TOKEN

WORKDIR /app

COPY . .
RUN go mod tidy

COPY go.mod go.sum ./
RUN go env -w GOPRIVATE=gitlab.humo.tj

RUN go mod download

RUN go build -o converterApi cmd/main.go

# Run stage

FROM hub.docker.humo.lab/nexus-repository-alpine

WORKDIR /app

COPY --from=builder /app/converterApi .

CMD ["./converterApi"]
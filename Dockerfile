FROM gitlab.humo.tj:5050/devops/nexus-repository/nexus-repository-alpine:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

RUN go build -o converterApi main.go

CMD ["./converterApi"]
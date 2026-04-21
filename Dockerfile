FROM hub.docker.humo.lab/nexus-repository-alpine
# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

RUN go build -o converterApi main.go

CMD ["./converterApi"]
FROM gitlab.humo.tj:5050/devops/nexus-repository/nexus-repository-alpine:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем скомпилированную программу и конфигурацию из стадии сборки.
COPY converterApi .

# Открываем порт для приложения
EXPOSE 8086

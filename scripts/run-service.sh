#!/bin/bash
# Сборка и запуск converterApi.
#
# Запускать можно откуда угодно: скрипт сам переходит в корень репозитория.
# Это обязательно - пути к config.json и транспортному ключу в коде
# относительные и считаются от рабочего каталога.
#
# Использование:
#   ./scripts/run-service.sh              собрать и запустить
#   ./scripts/run-service.sh --no-build   запустить уже собранный бинарник

set -eu

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

BINARY="./converterApi"

if [ "${1:-}" != "--no-build" ]; then
    # go.mod требует 1.25, и при GOTOOLCHAIN=auto go подтягивает нужную версию
    # сам. Но если в окружении жёстко задан GOROOT от другой версии (например,
    # export GOROOT=/usr/local/go в /etc/profile.d/go.sh), сборка падает на
    # каждом пакете с "compile: version ... does not match go tool version".
    # Снимаем переменную: современный go определяет GOROOT самостоятельно.
    unset GOROOT

    echo "Сборка ($(go version | awk '{print $3}'))..."
    go build -o "$BINARY" cmd/main.go
fi

if [ ! -x "$BINARY" ]; then
    echo "Бинарник $BINARY не найден - запустите без --no-build" >&2
    exit 1
fi

PORT=$(sed -n 's/.*"port"[[:space:]]*:[[:space:]]*"\([0-9]*\)".*/\1/p' internal/config/config.json)
PORT="${PORT:-8086}"

if command -v ss >/dev/null 2>&1 && ss -ltn 2>/dev/null | grep -q ":$PORT "; then
    echo "Порт $PORT уже занят - остановите работающий экземпляр" >&2
    exit 1
fi

echo "Запуск на порту $PORT (Ctrl+C для остановки)"
exec "$BINARY"

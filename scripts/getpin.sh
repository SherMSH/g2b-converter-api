#!/bin/bash
# Проверка PIN по карте через POST /g2b/VerifyPIN.
#
# Прочитать PIN нельзя ни через наш сервис, ни через D8: процессинг хранит его
# в виде PVV и наружу не отдаёт (спецификация, 7.7: "Due to security reasons PIN
# is not returned in the response"). Поэтому единственное, что можно сделать -
# проверить, подходит ли конкретное значение.
#
# ВНИМАНИЕ: неудачные попытки увеличивают счётчик неверных вводов в процессинге
# и в итоге блокируют карту. Сброс - операцией ResetBadPINTriesRq.
#
# Использование:
#   ./scripts/getpin.sh <PAN> <PIN> [срок YYMM]
#
# Настройки через переменные окружения:
#   HOST     адрес сервиса, по умолчанию http://localhost:8086
#   API_KEY  ключ, по умолчанию берётся из internal/config/config.json

set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

PAN="${1:-}"
PIN="${2:-}"
EXPDATE="${3:-3004}"

if [ -z "$PAN" ] || [ -z "$PIN" ]; then
    echo "Использование: $0 <PAN> <PIN> [срок YYMM]" >&2
    exit 1
fi

HOST="${HOST:-http://localhost:8086}"

if [ -z "${API_KEY:-}" ]; then
    API_KEY=$(sed -n 's/.*"api_key"[[:space:]]*:[[:space:]]*"\(.*\)".*/\1/p' \
        "$ROOT/internal/config/config.json" 2>/dev/null)
fi

if [ -z "$API_KEY" ]; then
    echo "Не найден API-ключ: задайте переменную API_KEY" >&2
    exit 2
fi

echo "--- POST $HOST/g2b/VerifyPIN  (PAN $PAN, срок $EXPDATE)"

# --noproxy обязателен: иначе даже localhost уходит в корпоративный прокси
curl -s -m 60 --noproxy '*' -X POST "$HOST/g2b/VerifyPIN" \
    -H "X-API-Key: $API_KEY" \
    -H 'Content-Type: application/json' \
    -d "{\"pan\":\"$PAN\",\"expiryDate\":\"$EXPDATE\",\"pin\":\"$PIN\"}" \
    -w '\n(HTTP %{http_code})\n'

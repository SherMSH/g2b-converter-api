# Deployment

> **Назначение:** как локально поднять проект, какие есть окружения, как работает CI/CD и деплой в продакшен.
> **Обновлено:** 2026-08

---

## Окружения

| Окружение | Назначение | Адрес | Доступ |
|-----------|-----------|-------|--------|
| `local` | Локальная разработка | `http://localhost:8086` | разработчик |
| `dev` | Ветка `dev`, образ `g2bdev:latest` | внутренний хост контура | команда |
| `production` | Ветка `main`, образ `g2bmain:latest` | `127.0.0.1:8086` на хосте, наружу через реверс-прокси | контур банка |

Отдельного staging нет: `dev` и `production` различаются только веткой, образом и стендом процессинга в `config.json`.

---

## Локальный запуск

### Требования
- Go `1.25.0` (см. `go.mod`)
- Доступ к стенду D8 (`processing.address`) — без него сервис поднимется, но все операции будут падать
- Docker — только если нужно собрать образ; для запуска не обязателен

### Шаги

```bash
# 1. Клонировать репозиторий
git clone git@gitlab.humo.tj:MustafokulovSh/g2b-converter-api.git
cd g2b-converter-api

# 2. Зависимости
go mod download

# 3. Проверить конфигурацию
#    internal/config/config.json уже в репозитории; поправить при необходимости
#    processing.address, app.port, app.debug_mode

# 4. Собрать
go build -o converterApi cmd/main.go

# 5. Запустить ИЗ КОРНЯ ПРОЕКТА
./converterApi
```

Ни БД, ни Redis, ни очередей поднимать не нужно.

После запуска:
- API: `http://localhost:8086`
- Health-check: `curl -H 'X-API-Key: <ключ>' http://localhost:8086/g2b/ping`
- Метрики: `http://localhost:8086/metrics`

### Частые проблемы

- **Паника `Can't setup configurations!` при старте** → приложение запущено не из корня проекта. Путь `internal/config/config.json` зашит в `cmd/main.go` и резолвится относительно рабочего каталога.
- **`500` на `POST /g2b/SetPIN`** → не найден `internal/app/files/transport_setpin.der`. Тот же относительный путь; в Docker-образ этот файл **не копируется**, его нужно монтировать томом.
- **`Portal connection failed` в логах каждые 30 секунд** → недоступен `processing.address`.
- **Порт занят** → `lsof -i :8086`; порт задаётся в `config.json`, а не флагом.
- **OFFLINE-файлы не обрабатываются** → проверить `jobs.conv_scanner.is_on`, права на каталог `in` и соответствие имени файла маскам `utils.OfflineReqTypes`.

---

## Конфигурация

Переменные окружения **не используются**. Единственный источник настроек — `internal/config/config.json`, читается один раз при старте; горячей перезагрузки нет.

| Параметр | Обязательно | Пример | Описание |
|----------|-------------|--------|----------|
| `app.server.host` | ✅ | `0.0.0.0` | адрес прослушивания |
| `app.server.port` | ✅ | `8086` | HTTP-порт |
| `app.server.name` | ⚪ | `Converter API` | попадает в поле `service` логов |
| `app.api_key` | ✅ | — | ключ для всех маршрутов `/g2b/*`, кроме `d8convert` |
| `app.debug_mode` | ⚪ | `true` | уровень логов Debug + подстановка тестового срока действия карты |
| `app.client_timeout_seconds` | ✅ | `30` | таймаут клиента портала D8 |
| `app.storage.basepath` | ✅ | `/app/files/...` | корень OFFLINE-файлов и каталог архива |
| `app.storage.in` / `.out` | ✅ | `/in`, `/out` | подкаталоги, обязаны начинаться со слэша |
| `processing.address` | ✅ | `http://d8-tprocweb1.humo.lab` | базовый URL процессинга, без завершающего слэша |
| `processing.extra.agreed_expdate_yymm` | ⚪ | `3004` | срок действия по умолчанию в debug-режиме |
| `jobs.conv_scanner.is_on` | ✅ | `true` | включение сканера OFFLINE-файлов |
| `jobs.conv_scanner.interval_seconds` | ✅ | `60` | период сканирования |

Объявлены в моделях, но не используются: `app.server.address`, `app.server.token`, `app.default_params`, `jobs.conv_scanner.query_limit`, вся секция `envelope`.

### Секреты

Текущее состояние — все секреты в репозитории:

| Секрет | Место |
|--------|-------|
| API-ключ партнёра | `internal/config/config.json` → `app.api_key` |
| Логин и пароль портала D8 | литерал в `pkg/d8-proc-web/auth.go` |
| Логины и SHA-256 паролей клерков | `auth.ClerksMap` в `internal/auth/models.go` |
| Транспортный RSA-ключ | `internal/app/files/transport_setpin.der` |

`internal/config/config.json` указан в `.gitignore`, но был добавлен в индекс раньше правила и продолжает отслеживаться. Просто убрать его из индекса нельзя: `Dockerfile` копирует файл в образ, сборка сломается.

Смена любого секрета сегодня означает: правка файла → пересборка образа → деплой.

---

## Сборка образа

`Dockerfile`, две стадии:

| Стадия | База | Действия |
|--------|------|----------|
| builder | `hub.docker.humo.lab/nexus-repository-golang-alpine3.23` | `go mod download`, `go mod tidy`, `go build -o converterApi cmd/main.go` |
| runtime | `hub.docker.humo.lab/nexus-repository-alpine` | копирование бинарника и `internal/config/config.json` |

В builder-стадии выставляется `GOPRIVATE=gitlab.humo.tj`.

> В образ попадает только `config.json`. Каталог `internal/app/files/` не копируется, поэтому в контейнере недоступны `transport_setpin.der` и `favicon.ico` — для работы `SetPIN` ключ нужно смонтировать томом.

---

## CI/CD

`.gitlab-ci.yml`, только ветки `main` и `dev`:

```
push в main/dev
      ↓
[build]  раннер gitlab-runner-test
      ↓  docker build -t $IMAGE:latest -t $IMAGE:$CI_COMMIT_SHORT_SHA
      ↓  docker login + docker push --all-tags
      ↓
[deploy] раннер g2b-runner
      ↓  docker compose pull
      ↓  docker compose up -d
```

Имя образа: `$CI_REGISTRY_IMAGE/g2b$CI_COMMIT_REF_SLUG` — у веток разные образы (`g2bmain`, `g2bdev`).

**Что стоит знать:**
- деплой запускается автоматически сразу после сборки, ручного апрува нет;
- тестов и линтеров в пайплайне нет — блокирует мерж только ревью;
- откат выполняется вручную: подставить тег `:$CI_COMMIT_SHORT_SHA` предыдущего образа в `docker-compose.yml` и выполнить `docker compose up -d`.

---

## Развёртывание

`docker-compose.yml`:

| Параметр | Значение |
|----------|----------|
| Контейнер | `g2bconverterapi` |
| Образ | `gitlab.humo.tj:5050/mustafokulovsh/g2b-converter-api/g2bmain:latest` |
| Порт | `127.0.0.1:8086 → 8086` |
| Перезапуск | `unless-stopped` |
| Команда | `sh -c './converterApi >> /app/logs/app.log 2>&1'` |
| Сеть | внешняя `default_net` |

Тома:

| Хост | Контейнер | Назначение |
|------|-----------|------------|
| `/srv/g2b/logs` | `/app/logs` | Логи |
| `/srv/g2b/files/history` | `…/OffLine` | Архив исходных пакетов |
| `/srv/g2b/files/in` | `…/OffLine/in` | Входящие пакеты |
| `/srv/g2b/files/out` | `…/OffLine/out` | Результаты |

### Чек-лист перед релизом в prod

- [ ] `go build ./...` и `go vet ./...` проходят
- [ ] Изменения проверены на стенде `dev`
- [ ] Если менялся `config.json` — проверено, что стенд процессинга указан верно
- [ ] Если менялись маппинги статусов или сумм — проверено на реальном запросе
- [ ] Известен тег предыдущего образа для отката

---

## Эксплуатация

### Остановка

По `SIGINT`/`SIGTERM` выполняется `server.Shutdown` с таймаутом **1 секунда** — запросы дольше секунды обрываются. Фоновые задачи gocron отдельно не останавливаются.

### Логи

- Формат: JSON (zerolog), поля `time`, `level`, `caller`, `service`.
- Уровень Debug включается флагом `app.debug_mode`.
- Файл: `/srv/g2b/logs/app.log` на хосте. Ротации нет — за размером файла надо следить.
- Префиксы контекста: `[MAIN]`, `[SERVER]`, `[JOBS]`, `[SERVICE]`.

> В логи пишутся тела запросов целиком: `SOAPLogger` — включая JSON `SetPIN` с открытым PIN, `ClerkAuth` (уровень Debug) — включая пароль клерка. Доступ к файлу логов приравнивается к доступу к учётным данным.

### Метрики

| Метрика | Тип | Метки |
|---------|-----|-------|
| `http_requests_total` | counter | `method`, `endpoint`, `status` |
| `http_request_duration_seconds` | histogram | `method`, `endpoint` |
| `http_errors_total` | counter | `method`, `endpoint`, `status` |
| `request_processing_seconds` | summary | `method`, `endpoint`, `status` |
| `service_status` | gauge | — (1/0 по последнему запросу) |
| `cpu_usage_percent`, `memory_usage_bytes`, `num_goroutines`, `service_uptime_seconds` | gauge/counter | — (обновляются раз в минуту) |

Метрик по исходящим вызовам процессинга нет — деградацию D8 видно только по времени ответа входящих запросов и по логам.

### Диагностика инцидентов

| Симптом | Что смотреть |
|---------|--------------|
| `401` на `/g2b/*` | `app.api_key` против заголовка `Authorization`/`X-API-Key` |
| `401 wrong clerk or password` | логин в `auth.ClerksMap`, SHA-256 пароля |
| SOAP Fault «Unknown XML body» | корневой тег тела не входит в `utils.BodyTypes` |
| Запросы висят ~90 секунд | недоступен `processing.address` (таймаут xapi) |
| `Portal connection failed` | недоступен портал D8 или сменились креды |
| Файлы копятся в `in/` | `jobs.conv_scanner.is_on`, права на тома, соответствие имени маскам |
| В `out/` лежит текст ошибки вместо результата | ошибка `Call()`; исходник искать в архиве `basepath` |

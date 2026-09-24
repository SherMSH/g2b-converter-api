# Architecture

> **Назначение:** описание архитектурного подхода, слоёв приложения, структуры папок и правил зависимостей между модулями.
> **Связано с:** [Stack.md](./Stack.md), [Database.md](./Database.md)
> **Обновлено:** 2026-08

---

## Архитектурный подход

Проект построен как **layered-адаптер (protocol gateway)**: единственная задача сервиса — преобразование протоколов, собственной бизнес-логики и собственных данных у него нет.

**Почему этот подход:**
- сервис не владеет состоянием — источник истины полностью на стороне процессинга D8, поэтому слой персистентности не нужен;
- набор операций диктуется внешним контрактом FIMI, а не внутренней моделью: удобнее держать по пакету на операцию, чем общую доменную модель;
- нагрузка на изменение приходится на маппинг полей, поэтому цена лишних абстракций выше их пользы.

**Следствия, которые надо учитывать:**
- слоя Domain в классическом смысле нет: «доменные» знания — это таблицы соответствия кодов в `internal/utils/compass.go`;
- каркас `repository → service → handlers` собирается в `app.New()`, но обработчики зарегистрированы как пакетные функции, а не как методы `Handler`, поэтому внедрение зависимостей фактически не работает — вся связность идёт через пакетные вызовы и глобальный `config.Config`.

---

## Слои

```
┌─────────────────────────────────────────────┐
│  Presentation: router + middlewares          │  ← маршруты, аутентификация, логирование, метрики
├─────────────────────────────────────────────┤
│  Handlers: D8Converter, SetPIN, conv-файлы   │  ← разбор запроса, диспетчеризация по типу
├─────────────────────────────────────────────┤
│  Models ONLINE/OFFLINE: request/response/Svc │  ← маппинг FIMI ↔ D8 по каждой операции
├─────────────────────────────────────────────┤
│  Service G2B                                 │  ← вызовы процессинга, сборка MDI
├─────────────────────────────────────────────┤
│  Infrastructure: utils.SendRequest,          │  ← HTTP-клиенты, файлы, crypto, logger, metrics
│  pkg/d8-proc-web, pkg/storage, pkg/crypto    │
└─────────────────────────────────────────────┘
```

### Правила зависимостей

- **Сверху вниз.** Handlers → models → service/G2B → pkg/*. Обратных вызовов нет.
- **Конфигурация — глобальная.** `config.Config` читается из любого слоя; это осознанное упрощение, но оно же делает пакеты нетестируемыми в изоляции.
- **`internal/utils` — общий низ.** Его импортируют все, сам он зависит только от `pkg/logger`.
- **Исключение:** `internal/middlewares` импортирует `internal/handlers` ради `SendSoapFault` — Presentation зависит от слоя обработчиков.

### Что можно и нельзя

| Откуда | Куда | Можно? |
|--------|------|--------|
| handlers | models/{ONLINE,OFFLINE} | ✅ |
| models | service/G2B | ✅ |
| service/G2B | pkg/* , internal/utils | ✅ |
| service/G2B | internal/handlers | ❌ |
| pkg/* | internal/models, internal/service | ❌ (исключение: `pkg/d8-proc-web` импортирует `D8CORP` ради `CommonResp`) |
| любой слой | `config.Config` | ✅ (по факту; в новом коде лучше передавать значения аргументами) |

---

## Структура папок

```
cmd/
└── main.go                  # init конфигов/логгера/клиентов, старт, graceful shutdown

internal/
├── app/                     # сборка Gin + http.Server
│   └── files/               # favicon, транспортный RSA-ключ transport_setpin.der
├── config/                  # config.json и его модели
├── router/                  # маршруты /g2b/*
├── middlewares/             # ClerkAuth, CheckApiKey, D8ProcWebAuth, CORS, Prometheus, SOAPLogger
├── handlers/                # D8Converter, SetPIN, PutConvFile, GetConvFile, SendSoapFault
├── auth/                    # разбор Clerk/Password из SOAP, карта клерков
├── models/
│   ├── request.go           # SoapEnvelope, интерфейсы MDIface / TrnInputIface
│   ├── ONLINE/<Op>/         # request.go + response.go + service.go на операцию
│   ├── OFFLINE/<Op>/        # пакетные XML-модели
│   ├── OFFLINE/models.go    # MRecord — универсальная запись пакета
│   └── D8CORP/              # JSON DTO процессинга: запросы, ответы, MDI
├── service/
│   ├── service.go           # Service + интерфейс G2bServiceIface
│   └── G2B/                 # клиент процессинга: карты, счета, клиенты, PIN, транзакции
├── jobs/                    # ConvScanner + переавторизация (gocron)
├── utils/                   # GetRqType, SendRequest, справочники кодов, константы
└── repository/              # заглушка, БД в проекте нет

pkg/
├── crypto/                  # 3DES/RSA, PIN-блоки (Format0, ZPK)
├── d8-proc-web/             # сессионный HTTP-клиент портала D8
├── logger/                  # zerolog
├── prometheus/              # метрики
└── storage/                 # чтение и перемещение OFFLINE-файлов

keys/                        # public.pem (резерв)
docs/                        # эта документация
```

Соглашение по ONLINE-операции — три файла в пакете:

| Файл | Содержимое |
|------|------------|
| `request.go` | структуры входящего XML + метод `Call()` — тонкая обёртка |
| `response.go` | структуры ответного `Envelope` |
| `service.go` | функция `Svc()` — вся логика: вызовы `service/G2B` и маппинг полей |

---

## Потоки данных

### ONLINE (синхронный SOAP)

```
Compass → POST /g2b/d8convert
        → ClerkAuth (Clerk + SHA-256 пароля)
        → D8Converter: проверка Content-Type, разбор конверта
        → utils.GetRqType(body) — тип по имени корневого тега
        → switch → models/ONLINE/<Op>.Call() → Svc()
        → service/G2B → HTTP к D8
        → сборка FIMI Envelope → c.XML(200)
```

Ключевое решение: тело конверта не разбирается на первом проходе — `XmlBody.XMLData` хранится как `innerxml`. Это позволяет держать 20+ несовместимых схем за одним эндпоинтом.

### OFFLINE (пакетные файлы)

```
Оператор → PUT /g2b/convFile → валидация MIME и имени → запись в storage.in
ConvScanner (gocron, 60 с) → glob по маскам → xml.Unmarshal → Root.Call()
        → service/G2B → MDI JSON в D8
        → архив исходника → запись результата → перенос в storage.out
Оператор → GET /g2b/convFile/:filename → файл из storage.out
```

Подробные пошаговые описания — в [Workflows/](./Workflows/).

---

## Онлайн-операции

Диспетчеризация — `switch` по `utils.RqBodyType` в `handlers/converter.go`.

| Операция | Обращается к процессингу |
|----------|--------------------------|
| `GetCardInfoRq` | `GetCardInfo` |
| `GetAcctInfoRq` | `GetAcctInfoG2b`, `GetCardsListG2b` |
| `GetAcctStatementRq` | `GetAcctInfoG2b`, `GetCardsListG2b`, `GetCardTransactionHistory` |
| `GetCardStatementRq` | `GetCardTransactionHistory` |
| `GetCVVRq` | `GetCVVG2b` |
| `GetPersonInfoRq` | `GetCustomerInfoG2b` |
| `GetTransInfoRq` | `GetCardInfo`, `GetTransactionDetailsG2b` |
| `POSRequestRq` | `GetTerminalByIdG2b`, `GetCardInfo`, `InitiateTransaction`, `AuthorizeTransaction`, `GetTransactionDetailsG2b` |
| `SetCardStatusRq` | `SetCardStatusG2b` |
| `SetAcctStatusRq` | `GetAcctInfoG2b`, `SetAccountStatusG2b` |
| `ResetBadPINTriesRq` | `ResetCardPINTriesG2b` |
| `VerifyPINRq` | `GetExpDateByPan`, `VerifyPinStatusG2b` |
| `SetDynamicPVV_PINOffsetRq` (только мобильный канал) | `GetExpDateByPan`, `SetPinG2b` |
| `UpdatePersonRq` | `UpdateCustomerG2b` |
| `UpdateCard2AcctLinkRq` | `UpdateCardAcctLinkG2b`, `DeleteCardAcctLinkG2b`, `SetCardStatusG2b` |
| `AddCMSAbonentRq`, `AddPersonCMSAbonentRq` | `AddCardNotificationG2b` |
| `ChangeCMSAbonentRq` | `AddCardNotificationG2b`, `DeleteCardNotificationG2b` |
| `RemoveCMSAbonentRq`, `RemovePersonCMSAbonentRq` | `DeleteCardNotificationG2b` |
| `AcctCreditRq`, `AcctDebitRq` | **заглушка** — фиксированные значения в `Svc()` |
| `InitSessionRq` | **заглушка** — вызов D8 закомментирован |

Особый случай — `UpdateCard2AcctLinkRq`: единственная операция, где повторно разбирается всё тело запроса, а не `Body.XMLData`, из-за собственной структуры конверта в пакете.

Зарезервированы в `utils.BodyTypes`, но не обрабатываются: `DeleteCard2AcctLinkRq`, `SetCardPersonRq` — попадают в `default` и получают SOAP Fault.

### Как добавить новую операцию

1. Константа в `utils.BodyTypes` (`internal/utils/constants.go`).
2. Пакет `internal/models/ONLINE/<Op>/` с тремя файлами.
3. Ветка `case` в `handlers/converter.go`.
4. При необходимости — функция в `internal/service/G2B/`.
5. Строка в таблице выше и, если процесс нетривиальный, файл в `Workflows/`.

---

## Сквозные практики

### Обработка ошибок
- Типизированных ошибок нет: используется `fmt.Errorf` с текстом.
- Наружу ONLINE-канал всегда отдаёт SOAP Fault через `handlers.SendSoapFault(c, status, code, text)`; файловые и PIN-эндпоинты — JSON.
- Текст ошибки сервиса подставляется в `Reason.Text` как есть — внутренние сообщения видны партнёру.

### Логирование
- Структурированный JSON (zerolog), уровень зависит от `app.debug_mode`.
- Префиксы контекста: `[MAIN]`, `[SERVER]`, `[JOBS]`, `[SERVICE]`.
- `correlationId` не пробрасывается — связать запись партнёра с вызовом D8 можно только по времени и содержимому.
- `SOAPLogger` пишет тело запроса целиком, включая PIN и пароль клерка. При работе с логами это надо учитывать.

### Транзакционность
- Транзакций в привычном смысле нет. Для финансовых операций используется цепочка D8: `initiateTransaction` → `authorizeTransaction` → при необходимости `reverseTransaction`.
- Идемпотентность обеспечивается `ecTxRefno`; для MDI-операций её нет — повторная обработка файла создаст дубли.

### Конкурентность
- `ConvScanner` защищён пакетным `sync.Mutex` — параллельных проходов не бывает.
- `d8procweb.Client` общий для всех горутин; cookie-сессия одна на процесс.
- `config.Config.Processing.Token` пишется фоновой задачей и читается обработчиками без синхронизации.

---

## Архитектурные решения

- **Без БД.** Сервис не хранит состояние: всё в D8 и в файлах. Это упрощает деплой, но лишает возможности повторить неудачную операцию и вести собственный аудит.
- **Один эндпоинт на все ONLINE-операции.** Соответствует контракту FIMI; цена — гигантский `switch` в `converter.go` (~380 строк) с повторяющимся кодом обработки ошибок.
- **Пакет на операцию.** Дублирование структур `Envelope`/`Response` между пакетами принято сознательно: схемы FIMI расходятся в деталях, общая структура получалась бы с десятками опциональных полей.
- **Два клиента к процессингу.** Исторически: прямые `xapi`-вызовы и портал `proc-web` с сессией. Единого клиента нет, таймауты и политика ошибок у них разные.

> Значимые изменения структуры стоит обсуждать до реализации: точек, где ошибка маппинга приводит к неверной финансовой операции, в проекте много.

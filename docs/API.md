# API

> **Назначение:** общие принципы построения API. Контракт ONLINE-канала задан XSD-схемами Compass FIMI, OpenAPI-спецификации у сервиса нет.
> **Обновлено:** 2026-08

---

## Базовая информация

- **Тип API:** SOAP 1.2 (ONLINE-канал) + REST/JSON и multipart (служебные операции).
- **Базовый URL:**
  - production: контейнер слушает `127.0.0.1:8086`, наружу — через реверс-прокси контура;
  - локально: `http://localhost:8086`.
- **Версионирование:** отсутствует. Версия протокола передаётся атрибутом `Ver` внутри запроса FIMI и копируется в ответ.
- **Спецификация:** XSD Compass FIMI (`fimi.xsd`, `fimi_types.xsd`) на стороне партнёра. Аннотации `@title`/`@version`/`@host` в `cmd/main.go` остались от неподключённого swag.

---

## Маршруты

| Метод | Путь | Аутентификация | Назначение |
|-------|------|----------------|------------|
| POST | `/g2b/d8convert` | Clerk | Все ONLINE-операции FIMI |
| GET | `/g2b/ping` | API-ключ | Health-check |
| PUT | `/g2b/convFile` | API-ключ | Загрузка OFFLINE-пакета |
| GET | `/g2b/convFile/:filename` | API-ключ | Выгрузка результата обработки |
| POST | `/g2b/SetPIN` | API-ключ | Установка PIN по карте |
| GET | `/metrics` | нет | Prometheus |
| GET | `/favicon.ico` | нет | Статика |
| POST | `/g2b/d8-proc-web/d8convert` | API-ключ + сессия портала | **обработчик не зарегистрирован** — маршрут возвращает пустой `200` |

Глобальные middleware: `CORS`, `Prometheus`. Для группы с API-ключом дополнительно `SOAPLogger`.

---

## Аутентификация

В сервисе **две независимые схемы** — их нельзя использовать взаимозаменяемо.

### Clerk — только `POST /g2b/d8convert`

Учётные данные передаются внутри тела запроса:

```xml
<Request Ver="1.0" Product="..." Echo="..." Session="..."
         Clerk="onlinebankapp" Password="<пароль>"/>
```

Проверка: логин ищется в `auth.ClerksMap`, пароль сравнивается как `hex(sha256(password))`. Разбор выполняется нестрогим XML-декодером, промежуточные теги игнорируются.

Ошибки: пустое тело → `400`, ошибка разбора → `500`, неверная пара → `401 wrong clerk or password`. Все — в формате SOAP Fault.

### API-ключ — остальные маршруты `/g2b/*`

```
Authorization: Bearer <api_key>
# или
X-API-Key: <api_key>
```

Префикс `Bearer ` отбрасывается с обеих сторон сравнения. Значение — `app.api_key` из `config.json`, одно на всех потребителей; ротация требует пересборки образа.

### Сессия портала D8 — `/g2b/d8-proc-web/*`

Middleware `D8ProcWebAuth` выполняет `Signin()` до обработчика и `Signout()` после. Это аутентификация **сервиса перед процессингом**, а не клиента перед сервисом.

---

## Формат запросов и ответов

### ONLINE: запрос

`Content-Type` обязан быть одним из: `application/soap+xml`, `text/xml`, `application/xml`. Иначе — `415` простым текстом.

```xml
<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <m:GetCardInfoRq xmlns:m="http://schemas.compassplus.com/two/1.0/fimi.xsd">
      <Request Ver="1.0" Product="..." Echo="..." Session="..."
               Clerk="..." Password="...">
        <PAN>...</PAN>
        <ExpirationDate>3004</ExpirationDate>
      </Request>
    </m:GetCardInfoRq>
  </s:Body>
</s:Envelope>
```

Два обязательных свойства структуры:

- **тип операции** определяется по имени корневого тега тела (`GetCardInfoRq`), а не по URL или параметру;
- **все поля payload — дочерние элементы `<Request>`**, а не соседние с ним. Если положить их рядом, `auth.ParseAuth` упадёт с `expected element type <Request> but have <PAN>` и запрос получит Fault `500 Auth error` ещё до диспетчеризации.

### ONLINE: успешный ответ

`Content-Type: application/xml; charset=utf-8`, HTTP `200`. Конверт всегда содержит три пространства имён:

| Атрибут | Значение |
|---------|----------|
| `xmlns:s` | `http://www.w3.org/2003/05/soap-envelope` |
| `xmlns:m0` | `http://schemas.compassplus.com/two/1.0/fimi_types.xsd` |
| `xmlns:m1` | `http://schemas.compassplus.com/two/1.0/fimi.xsd` |

Общие поля `Response`: `Ver`, `Product`, `Echo` копируются из запроса, `ResponseAttr="1"`, `TranId` генерируется из времени.

### ONLINE: ошибка

`Content-Type: application/soap+xml; charset=utf-8`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Fault>
  <Code>
    <Value>soap:Client</Value>
  </Code>
  <Reason>
    <Text>Service error: ...</Text>
  </Reason>
</Fault>
```

> Текст ошибки внутреннего сервиса подставляется в `Reason.Text` как есть — сообщения процессинга видны партнёру.

### REST-эндпоинты

| Эндпоинт | Запрос | Успешный ответ |
|----------|--------|----------------|
| `GET /g2b/ping` | — | `{"message":"Pong!"}` |
| `POST /g2b/SetPIN` | `{"pan":"...","expiryDate":"...","pin":"..."}` | `200` **с пустым телом** |
| `PUT /g2b/convFile` | multipart, поле `file` | `{"message":"Файл успешно загружен","filename":"...","size":N}` |
| `GET /g2b/convFile/:filename` | — | поток файла, `application/octet-stream` |

Ошибки REST-эндпоинтов — JSON вида `{"error":"описание"}`, сообщения на русском языке.

---

## HTTP-коды

| Код | Когда используем |
|-----|------------------|
| 200 | Успех (в том числе SOAP-ответ и пустой ответ `SetPIN`) |
| 204 | Preflight `OPTIONS` |
| 400 | Ошибка разбора тела, некорректное имя файла, неподходящий тип файла |
| 401 | Нет ключа / неверный ключ / неверная пара клерк-пароль |
| 404 | Файл не найден в каталоге `out` |
| 415 | Неподдерживаемый `Content-Type` для ONLINE-запроса |
| 500 | Ошибка разбора XML, ошибка файловой системы, ошибка вызова процессинга |

Особенность: ошибка бизнес-уровня в ONLINE-канале отдаётся как `400` с SOAP Fault, а не как `422`.

---

## Коды ошибок

Собственного реестра кодов у сервиса нет. Наружу транслируются:

| Источник | Поле | Значения |
|----------|------|----------|
| SOAP Fault | `Code.Value` | `soap:Client` — единственное используемое значение |
| Процессинг D8 | `status.rspcode` | `00` — успех; остальное — ошибка (текст в `status.message`) |
| Авторизация D8 | `D8RspCode` | `00`, `01`, `02`, `80`, `81`, `82`, `85` — см. [Glossary.md](./Glossary.md) |
| MDI | `C_ACTIONCODE`, `C_RSPCODE`, `I_REJMSG` | `0` — успех; иначе текст отказа |

---

## Соглашения

- **Имена маршрутов** повторяют имена операций FIMI (`/SetPIN`, `/convFile`) — camelCase и PascalCase вперемешку, единого стиля нет.
- **Идентификаторы** в ответах генерируются из времени (`utils.GenerateTimestampID`), UUID не используются.
- **Пагинация** есть только у истории транзакций D8 (`PagingParams`), наружу не выставлена.
- **Rate limiting** не реализован.
- **Идемпотентность** обеспечивается только для финансовых операций через `ecTxRefno`; повторная загрузка OFFLINE-файла приведёт к повторной обработке.

---

## Ограничения

| Ограничение | Значение |
|-------------|----------|
| Размер multipart в памяти | 10 MiB (`router.MaxMultipartMemory`) |
| Таймаут вызова D8 xapi | 90 секунд |
| Таймаут вызова портала D8 | `app.client_timeout_seconds` (30 секунд) |
| Таймаут graceful shutdown | 1 секунда — долгие запросы обрываются |
| CORS | `Access-Control-Allow-Origin: *`, методы `GET, POST, PUT, OPTIONS` |

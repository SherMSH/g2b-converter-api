# SetDynamicPVV_PINOffsetRq — установка PIN из мобильного приложения

**Куда:** `POST /g2b/mobile/v1/d8convert` — только мобильный канал.
В основном `/g2b/d8convert` операции нет, там придёт `Unknown XML body`.

**Заголовки:**

```
Content-Type: application/soap+xml     (допустимы также text/xml и application/xml)
X-API-Key:    <ключ мобильного канала>
```

Учётные данные клерка передаются атрибутами `Clerk` и `Password` в `<fimi:Request>`.

---

## Запрос

```xml
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope"
               xmlns:fimi="http://schemas.compassplus.com/two/1.0/fimi.xsd"
               xmlns:fimi1="http://schemas.compassplus.com/two/1.0/fimi_types.xsd">
   <soap:Header/>
   <soap:Body>
      <fimi:SetDynamicPVV_PINOffsetRq>
         <fimi:Request Ver="3.5" Product="FIMI" Echo="7" Clerk="..." Password="...">

            <!-- Номер карты. Обязательное поле -->
            <fimi1:PAN>5058270530003879</fimi1:PAN>

            <!-- Доп. номер карты. Не используется, можно не присылать -->
            <fimi1:MBR>0</fimi1:MBR>

            <!-- Не используется -->
            <fimi1:CardUID></fimi1:CardUID>

            <!-- ДОЛЖЕН БЫТЬ ПУСТЫМ. PIN выбирает процессинг и наружу не отдаёт
                 (спецификация D8, 7.7), задать своё значение нельзя. Новый PIN
                 доходит до держателя SMS-оповещением, поэтому у карты должен
                 быть заведён контракт SMSGEN.
                 Присланное значение = отказ, а не молчаливое игнорирование -->
            <fimi1:PINBlock></fimi1:PINBlock>

            <!-- Не используется: рабочих ключей в схеме процессинга нет -->
            <fimi1:KeyId></fimi1:KeyId>

            <!-- Принимается, но не применяется: процессинг индекс ключа не меняет -->
            <fimi1:PVKI>0</fimi1:PVKI>

            <!-- Срок действия YYMM. Необязательно: если не прислать,
                 сервис возьмёт срок из данных карты -->
            <fimi1:ExpDate>3004</fimi1:ExpDate>

            <!-- НЕ ДЕЙСТВУЕТ. Одноразового PIN в процессинге нет, PIN действует
                 постоянно. Значение 1 не отклоняется, но пишется
                 предупреждение в лог -->
            <fimi1:SingleOperation>1</fimi1:SingleOperation>

            <!-- Произвольная причина изменения -->
            <fimi1:ChangeReason>123</fimi1:ChangeReason>

         </fimi:Request>
      </fimi:SetDynamicPVV_PINOffsetRq>
   </soap:Body>
</soap:Envelope>
```

## Успешный ответ

HTTP `200`, `Content-Type: application/xml; charset=utf-8`.

```xml
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope"
            xmlns:m1="http://schemas.compassplus.com/two/1.0/fimi.xsd"
            xmlns:m0="http://schemas.compassplus.com/two/1.0/fimi_types.xsd">
   <s:Body>
      <m1:SetDynamicPVV_PINOffsetRp>

         <!-- Echo       - возвращается из запроса как есть. По нему приложение
                           сопоставляет ответ со своим запросом
              Product    - тоже из запроса
              Response   - признак успеха. "1" = PIN сгенерирован.
                           Отказ приходит не здесь, а SOAP Fault'ом (см. ниже)
              TranId     - идентификатор обращения на нашей стороне,
                           timestamp-based. Пригодится при разборе инцидентов:
                           по нему мы найдём запрос в логе
              Ver        - версия формата ответа, всегда "1.0" -->

         <m1:Response Echo="42" Product="FIMI" Response="1"
                      TranId="1790259899756691687" Ver="1.0"/>

      </m1:SetDynamicPVV_PINOffsetRp>
   </s:Body>
</s:Envelope>
```

Тело ответа пустое: **значения PIN в нём нет и не будет**. Успех означает лишь
то, что процессинг принял команду и сгенерировал новый PIN.

Комментарии выше добавлены для пояснения - в реальном ответе их нет.

## Отказ

Приходит SOAP Fault с HTTP `400`. Причина - в `<Text>`, префикс
`Service error:` добавляет наш сервис.

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Fault>
  <Code>
    <Value>soap:Client</Value>          <!-- всегда soap:Client -->
  </Code>
  <Reason>
    <!-- Текст причины. Сюда попадает либо наша проверка, либо ответ
         процессинга в формате "<rspcode> - <сообщение>" -->
    <Text>Service error: 00 - Card not found</Text>
  </Reason>
</Fault>
```

## Откуда держатель узнаёт PIN

Значение не возвращается никому — ни вызывающей стороне, ни нам. Процессинг
отправляет его держателю карты сам, оповещением. Поэтому карта должна иметь
контракт `SMSGEN` с актуальным номером, иначе PIN сменится, а сообщить его
будет некому.

## Проверка результата

Тем же каналом PIN не проверяется — для этого есть `VerifyPINRq` в основном
`/g2b/d8convert` или JSON-маршрут `POST /g2b/VerifyPIN`:

```
POST /g2b/VerifyPIN
X-API-Key: <ключ>
{"pan":"5058270530003879","expiryDate":"3004","pin":"1234"}

→ {"verified": false, "reason": "bad PIN"}
```

Учтите: **неудачные проверки увеличивают счётчик неверных вводов**, лимит — 3.
Сброс — операцией `ResetBadPINTriesRq`.

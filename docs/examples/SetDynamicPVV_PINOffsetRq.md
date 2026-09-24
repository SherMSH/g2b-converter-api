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

            <!-- ВНИМАНИЕ: здесь ОТКРЫТЫЙ PIN, 4-12 цифр, а не блок под TPK.
                 PIN-блок под рабочим ключом принять нельзя: этого ключа у нас
                 нет. Процессинг требует блок под одноразовым 3DES-ключом,
                 который сервис собирает сам и шифрует RSA-ключом процессинга.
                 Пустое значение = отказ: сбрасывать в D8 нечего -->
            <fimi1:PINBlock>5810</fimi1:PINBlock>

            <!-- Не используется: рабочих ключей в схеме процессинга нет -->
            <fimi1:KeyId></fimi1:KeyId>

            <!-- Принимается, но не применяется: setPIN индекс ключа не меняет -->
            <fimi1:PVKI>0</fimi1:PVKI>

            <!-- Срок действия YYMM. Необязательно: если не прислать,
                 сервис возьмёт срок из данных карты -->
            <fimi1:ExpDate>3004</fimi1:ExpDate>

            <!-- НЕ ДЕЙСТВУЕТ. Одноразового PIN в процессинге нет, PIN меняется
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
         <m1:Response Echo="7" Product="FIMI" Response="1" TranId="1790257728046629335" Ver="1.0"/>
      </m1:SetDynamicPVV_PINOffsetRp>
   </s:Body>
</s:Envelope>
```

`Echo` и `Product` возвращаются из запроса, `TranId` — идентификатор обращения на нашей стороне.

## Отказы

Приходят SOAP Fault с HTTP `400`, текст в `<Text>`:

| Ситуация | Текст |
|----------|-------|
| пустой `PAN` | ``wrong mandatory field `fimi1:PAN` `` |
| пустой `PINBlock` | `PINBlock is empty: сброс динамического PVV процессингом не поддерживается` |
| в `PINBlock` шифрованный блок, а не PIN | `PINBlock must contain clear PIN (4-12 digits): PIN-блок под рабочим ключом процессинг не принимает` |
| карта не найдена, PIN не установлен | текст ошибки процессинга, например `09 - System malfunction` |

## Проверка результата

Тем же каналом PIN не проверяется — для этого есть `VerifyPINRq` в основном
`/g2b/d8convert` или JSON-маршрут `POST /g2b/VerifyPIN`:

```
POST /g2b/VerifyPIN
X-API-Key: <ключ>
{"pan":"5058270530003879","expiryDate":"3004","pin":"5810"}

→ {"verified": true}
```

Учтите: **неудачные проверки увеличивают счётчик неверных вводов**, лимит — 3.
Сброс — операцией `ResetBadPINTriesRq`.

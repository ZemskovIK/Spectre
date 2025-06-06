# Сервис писем военных лет

## Описание проекта

Данный модуль представляет собой REST API CRUD-сервис. Он предназначен для работы с письмами военных лет, позволяя создавать, читать, обновлять и удалять записи этих исторических документов. Все данные писем передаются в зашифрованном виде.

## Флоу работы сервиса

Взаимодействие с сервисом строится по следующей схеме:

### Получение писем

1. **(client) → [spectre] → [db]**  
   Клиент отправляет запрос на получение писем в сервис `spectre`, который обращается к базе данных.

2. **[db] → [spectre] → [crypto]**  
   База данных возвращает письма сервису, который затем отправляет их на внешний криптосервис для шифрования.

3. **[crypto] → [spectre] → (client)**  
   Криптосервис возвращает зашифрованные данные, которые сервис писем отправляет клиенту.

---

### Добавление и обновление писем

1. **(client) → [spectre] → [crypto]**  
   Клиент отправляет данные письма в сервис, который сразу пересылает их на криптосервис для шифрования.

2. **[crypto] → [spectre] → [db]**  
   После получения зашифрованных данных сервис сохраняет их в базе данных.

3. **[spectre] → (client)**  
   Сервис подтверждает успешное выполнение операции клиенту.

---

### Удаление писем

Удаление писем не требует шифрования.  
Сервис проверяет права пользователя и удаляет запись из базы данных, не взаимодействуя с криптосервисом.

---

## Функциональность

- **Авторизация**: Сервис работает с уровнями доступа клиента, админа.
- **Создание**: Добавление новых писем и пользователей в систему.
- **Чтение**: Получение и просмотр сохранённых писем и пользователей.
- **Обновление**: Изменение существующих записей писем.
- **Удаление**: Удаление писем и пользователей из системы.

---

## Структура модуля

- `cmd/` - Точка входа в приложение
    - `migrator/` - Приложение для создания или отката миграций
    - `spectre/` - Основное приложение
- `internal/` - Основная бизнес-логика и сервисы
    - `models/` - Модели приложения
    - `storage/` - Работа с файловой системой (базами данных)
        - `sqlite/` - Одна из реализаций базы (sqlite)
    - `srv/` - Работа с http
        - `proxy/` - Работа с обращением к внешнему сервису `crypto`
        - `lib/` - Вспомогательные функции широкого назначения
            - `methods/` - Функции для унификации методов в роутере
            - `response/` - Описание ответа
        - `api/` - Все что касается самого сервиса с письмами
            - `handlers/` - api обработчики запросов
        - `auth/` - Пакет авторизации (JWT)
- `pkg/` - Общие утилиты и вспомогательные функции
    - `logger/` - Обертка вокруг логгера (logrus)

---

## Обмен ключами с помощью Диффи-Хелмана

### Получить публичный ключ сервера K

`GET /ecdh`
С пустым телом и без токена.

**Успех (200 OK):**
```json
{
  "key": "base64string"
}
```

**Ошибка (502):**
```json
{
  "error": "error!",  
}
```

---

### Отправить свой ключ серверу

`POST /ecdh`
```json
{
  "key": "base64string"
}
```
**Успех (204 OK):**
Пустой ответ.

**Ошибка (502):**
```json
{
  "error": "error!",  
}
```

---

## Авторизация с JWT

`POST /login`
```json
{
  "login": "login",
  "password": "pass"
}
```

**Успех (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTA1LTIzVDE5OjM1OjIxLjkzOTA4NTY0NiswMzowMCIsInJvbGUiOjAsInN1YiI6MH0.B_WbK2YI13Cv3yju1FuQ6cGU1kE1G1_nAE8UYDgMuZA"
}
```

---


## Формат ответа API

Все ответы API имеют следующую структуру:

```python
{
  "content": { ... },   // полезная нагрузка (зашифрованные данные)
  "error": null | string, // описание ошибки, если она есть
  "iv": "base64",       // IV (инициализационный вектор) для расшифровки
  "hmac": "base64",     // HMAC для проверки целостности
  "nonce": "base64"     // nonce, если используется
}
```

- Если запрос успешен, поле `error` будет `null`.
- Все данные в поле `content` передаются в зашифрованном виде (например, как base64-строка или объект с полем `cipher_bytes`).

---

## API эндпоинты (в работе)

## Письма

### Получить все письма

`GET /api/letters`

**Успех (200 OK):**
```json
{
    "content": "base64string...TCm0oI...QDCCFs=",
    "error": null,
    "iv": "base64string...TCm0oI...QDCCFs=",
    "hmac": "base64string...TCm0oI...QDCCFs=",
    "nonce": "base64string...TCm0oI...QDCCFs="
}
```

**Ошибка (500 | 502):**
```json
{
  "content": null,
  "error": "error!",
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

---

### Получить письмо по ID

`GET /api/letters/{letter_id}`

**Успех (200 OK):**
```json
{
    "content": "base64string...TCm0oI...QDCCFs=",
    "error": null,
    "iv": "base64string...TCm0oI...QDCCFs=",
    "hmac": "base64string...TCm0oI...QDCCFs=",
    "nonce": "base64string...TCm0oI...QDCCFs="
}
```

**Ошибка (500 | 502 | 400 | 404 | 403):**
```json
{
  "content": null,
  "error": "error!",
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

---

### Удалить письмо по ID

`DELETE /api/letters/{letter_id}`

**Успех (200 OK):**
```json
{
  "content": null,
  "error": null,
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

**Ошибка (500 | 403 | 400 | 404):**
```json
{
  "content": null,
  "error": "error!",
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

---

### Создать новое письмо

`POST /api/letters`

**Тело запроса (application/json):**
```json
{
    "content": "base64string...TCm0oI...QDCCFs=",
    "error": null,
    "iv": "base64string...TCm0oI...QDCCFs=",
    "hmac": "base64string...TCm0oI...QDCCFs=",
    "nonce": "base64string...TCm0oI...QDCCFs="
}
```

**Успех (200 OK):**
```json
{
  "content": null,
  "error": null,
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

**Ошибка (500 | 502 | 403 | 400 | 422):**
```json
{
  "content": null,
  "error": "error!",
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

---

### Обновить письмо по ID

`PUT /api/letters/{letter_id}`

**Тело запроса (application/json):**
```json
{
    "content": "base64string...TCm0oI...QDCCFs=",
    "error": null,
    "iv": "base64string...TCm0oI...QDCCFs=",
    "hmac": "base64string...TCm0oI...QDCCFs=",
    "nonce": "base64string...TCm0oI...QDCCFs="
}
```

**Успех (200 OK):**
```json
{
  "content": null,
  "error": null,
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

**Ошибка (500 | 502 | 403 | 400 | 404):**
```json
{
  "content": null,
  "error": "error!",
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

---

## Пользователи (админ права)

### Получить пользователя по ID

`GET /api/users/{user_id}`

**Успех (200 OK):**
```json
{
    "content": {
        "id": 1,
        "login": "kek",
        "pass_hash": "YXNk",
        "access_level": 3
    },
    "error": null,
    "iv": "",
    "hmac": "",
    "nonce": ""
}
```

**Ошибка (500 Internal Server Error):**
```json
{
  "content": null,
  "error": "error!",
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

---

### Удалить пользователя по ID

`DELETE /api/users/{user_id}`

**Успех (200 OK):**
```json
{
  "content": null,
  "error": null,
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

**Ошибка (500 Internal Server Error):**
```json
{
  "content": null,
  "error": "error!",
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

---

### Создать нового пользователя

`POST /api/users`

**Тело запроса (application/json):**
```json
{
    "login":"asd1",
    "password":"asd1_pass",
    "access_level":2
}
```

**Успех (200 OK):**
```json
{
  "content": null,
  "error": null,
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

**Ошибка (422 | 400):**
```json
{
  "content": null,
  "error": "error!",
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

---

`DELETE /api/users/{user_id}`

**Успех (200 OK):**
```json
{
  "content": null,
  "error": null,
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

**Ошибка (500 Internal Server Error):**
```json
{
  "content": null,
  "error": "error!",
  "iv": "",
  "hmac": "",
  "nonce": ""
}
```

---

## Примечания

- Все поля, связанные с шифрованием (`content`, `iv`, `hmac`, `nonce`), обязательны для корректной работы клиента.
- Клиент отвечает за расшифровку данных и проверку целостности.
- Для доступа к API требуется JWT-токен, который передаётся в заголовке `Authorization: Bearer <TOKEN>`.
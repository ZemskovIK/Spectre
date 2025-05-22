# Сервис писем военных лет

## Описание проекта

Данный модуль представляет собой REST API CRUD-сервис. Он предназначен для работы с письмами военных лет, позволяя создавать, читать, обновлять и удалять записи этих исторических документов. Все данные писем передаются в зашифрованном виде.

---

## Возможности

- **Создание**: Добавление новых писем в систему.
- **Чтение**: Получение и просмотр сохранённых писем.
- **Обновление**: Изменение существующих записей писем.
- **Удаление**: Удаление писем из системы.
- **Авторизация**: Доступ к письмам осуществляется по JWT-токену, выдаваемому после логина.

---

## Структура модуля

- `cmd/` - Точка входа в приложение
    - `migrator/` - Приложение для создания или отката миграций
    - `spectre/` - Основное приложение
- `internal/` - Основная бизнес-логика и сервисы
    - `models/` - Модели приложения
    - `lib/` - Вспомогательные функции широкого назначения
    - `storage/` - Работа с файловой системой (базами данных)
        - `sqlite/` - Одна из реализаций базы (sqlite)
    - `server/` - Работа с http
        - `api/` - Все что касается самого сервиса с письмами
        - `auth/` - Пакет авторизации (JWT)
        - `methods/` - Функции для унификации методов в роутере
        - `response/` - Описание ответа
- `pkg/` - Общие утилиты и вспомогательные функции
    - `logger/` - Обертка вокруг логгера (logrus)

---

## Авторизация с JWT

`POST /login`
```json
{
  "login": "login1",
  "password": "pass1"
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

## API эндпоинты

### Получить все письма

`GET /api/letters`

**Успех (200 OK):**
```json
{
  "content": [
    "base64string1",
    "base64string2",
    "base64string3"
  ],
  "error": null,
  "iv": "base64string",
  "hmac": "base64string",
  "nonce": "base64string"
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

### Получить письмо по ID

`GET /api/letters/{letter_id}`

**Успех (200 OK):**
```json
{
  "content": "base64string",
  "error": null,
  "iv": "base64string",
  "hmac": "base64string",
  "nonce": "base64string"
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

### Создать новое письмо

`POST /api/letters`

**Тело запроса (application/json):**
```json
{
  "cipher_bytes": "base64string"
}
```
где шифруется структура вида
```json
{
  "author": "Author Name",
  "found_at": "2005-11-11T00:00:00Z",
  "found_in": "Место находки",
  "body": "Текст письма"
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

### Обновить письмо по ID

`PUT /api/letters/{letter_id}`

**Тело запроса (application/json):**
```json
{
  "cipher_bytes": "base64string"
}
```
где шифруется структура вида
```json
{
  "author": "Author Name",
  "found_at": "2005-11-11T00:00:00Z",
  "found_in": "Место находки",
  "body": "Текст письма"
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

**Ошибка (500 | 400 | 404):**
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

- Все поля, связанные с шифрованием (`cipher_bytes`, `iv`, `hmac`, `nonce`), обязательны для корректной работы клиента.
- Клиент отвечает за расшифровку данных и проверку целостности.
- Для доступа к API требуется JWT-токен, который передаётся в заголовке `Authorization: Bearer <TOKEN>`.
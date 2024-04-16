## Avito-test-assignment

[Задача](https://github.com/avito-tech/backend-trainee-assignment-2024?tab=readme-ov-file)

## Запуск
1. Создаем в папке configs файл .env с содержимым файла example.env
2. Пишем в консоль docker compose up --build

## Примеры

### Пример 1: Успешный запрос

**Запрос:**
```bash
curl -X GET "http://localhost:8088/banner"
     -H "token: 123"
```

**Ответ: Status 200**
```json
{
  "BannerID": 1,
  "FeatureID": 1,
  "Title": "Баннер 3",
  "Text": "Текст 1",
  "URL": "Урл 1",
  "IsActive": true,
  "CreatedAt": "2024-04-15T02:17:16.253917Z",
  "UpdatedAt": "2024-04-15T02:17:16.253917Z",
  "TagIDs": [1, 2, 3]
}
```

### Пример 2: Ошибка аутентификации

**Запрос:**
```bash
curl -X GET "http://localhost:8088/banner"
     -H "token:"
```

**Ответ: Status 401**
```
Unauthorized
```

### Пример 3: Ошибка доступа

**Запрос:**
```bash
curl -X GET "http://localhost:8088/banner"
     -H "token: 321"
```

**Ответ: Status 403**
```
Forbidden
```

### Пример 4: Ошибка данных

**Запрос:**
```bash
curl -X POST "http://localhost:8088/banner"
     -H "token: 321"
     -H "Content-Type: application/json"
     -d '{
         "content": {
             "title": "Баннер 3"
         }
     }'
```

**Ответ: Status 400**
```
Insufficient data to create a banner
```
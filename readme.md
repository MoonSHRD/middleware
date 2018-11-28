# Middleware
Сервис предназначен для авторизации клиентов.

# Способ взаимодействия
На данный момент выдает JWT токены для последующего взаимодействия с центрифугой.

При подключении вы получаете хэш, который надо подписать приватным ключом
(алгоритм подписи идентичен Ethereum) и отправить обратно,
если все правильно, то вернет JWT токен, который истечет через 5 минут.
Для получения нового токена - достаточно отправить любые данные.

## Соединение
Сервис поддерживает подключение только по websocket

## Тип данных
Сервис принимает и оправляет данные в виде json.

### Входящие:
```
type InData struct {
    Address string `json:"address"`
    Type    string `json:"type"`
    Data    string `json:"data"`
}
```

### Исходящие:
```
type OutData struct {
    Type    string `json:"type"`
    Data    string `json:"data"`
}
```

## Пример
1. Подключаемся к сервису и получаем следующий JSON:
    ```
    {
      "type": "auth",
      "data": {
        "hash": "742bc42264f857dc68331cc5c26d0f89474fb499a17ac35d8c84cf8491906b54"
      }
    }
    ```
2. Подписываем и отправляем:
    ```
    {
      "type": "auth",
      "data": "0xc6f048278a26b2c27e72ea02476670818021f9f13c609e27473339c11015ae8e6e0c945efbbdc0415f88b030f1d4ec6b8d454c68543cd85f3ab33a92a8689b8b1b",
      "address": "0x2d5e26a23c5190d88ee424e6a180f1fc8fd62ec5"
    }
    ```
    data - наша подпись
3. Получаем JWT
    ```
    {
      "type": "auth_success",
      "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDM0MTA4NDUsInN1YiI6IjB4MmQ1ZTI2YTIzYzUxOTBkODhlZTQyNGU2YTE4MGYxZmM4ZmQ2MmVjNSJ9.Ca7RG4IgG7ldFYqwoJJKJcZ34IDLwXbvGykCGihs1tY"
    }
    ```
4. При истечении токена отправляем любые данные
    ```
    {
      "type": "fafafa",
      "data": "fwafwa",
      "address": "0x2d5e26a23c5190d88ee424e6a180f1fc8fd62ec5"
    }
    ```
5. Полчаем новый JWT
    ```
    {
      "type": "newJWT",
      "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDM0MTEwODksInN1YiI6IjB4MmQ1ZTI2YTIzYzUxOTBkODhlZTQyNGU2YTE4MGYxZmM4ZmQ2MmVjNSJ9.Q7Jy7FJw_07UUuHcPMs9jd1rurZQLmrswiiGr2V4py8"
    }
    ```
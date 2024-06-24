# goauth
Сервис для аутентификации, написанный на Go

## Установка и конфигурация
- Склонировать репозиторий:
  ```
  git clone https://github.com/mx4alex/goauth.git
  ```
- Настроить конфигурацию в файле `config.yaml` и окружение в файле `.env` 
- Запустить *docker compose*
  ```
  docker compose up --build
  ```

## Использование

### Сервис поддерживает следующие эндпоинты:
- `POST /auth/sign-up` регистрирует нового пользователя
- `POST /auth/sign-in` авторизирует пользователя
- `POST /auth/refresh` обновляет refresh token пользователя

Документация находится в папке <a href="https://github.com/mx4alex/goauth/tree/main/docs">docs</a>

Визуальная документация Swagger UI доступна по адресу [`http://localhost:8080/swagger/index.html#`](http://localhost:8080/swagger/index.html#)

### Формат запросов

#### Sign Up
* Метод: `POST`
* Эндпоинт: `http://localhost:8080/auth/sign-up`
* Формат запроса:
```json
{
    "name": "ivan",
    "username": "ivan007",
    "password": "qwerty12345"
}
```
* Формат ответа:
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTM1MDUxNzIuOTc1MzgzLCJpYXQiOjE2OTM1MDQyNzIuOTc1NDEyLCJ1c2VybmFtZSI6Iml2YW4wMDcifQ.5ZB7_QHbohRxxhbxtuBqOGwvO-bZ2zoD9g5jLF9O9Zk",
    "refresh_token": "3423e6672e811489b97e99640120aa11cf74e8dbddd7ee6e75742a7abef74066"
}
```

#### Sign In

* Метод: `POST`
* Эндпоинт: `http://localhost:8080/auth/sign-in`
* Формат запроса:
```json
{
    "username": "ivan007",
    "password": "qwerty12345"
}
```
* Формат ответа:
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTM1MDUyNDEuMjAyMDE2LCJpYXQiOjE2OTM1MDQzNDEuMjAyMDE4LCJ1c2VybmFtZSI6Iml2YW4wMDcifQ.gcxo4MlHCkgUT9QGp--73rmBJL7FENqbScBx-qEcNOY",
    "refresh_token": "188ea22713b990a71f66f664534f8670151a5846f32ac04641110988d147610c"
}
```

#### Refresh

* Метод: `POST`
* Эндпоинт: `http://localhost:8080/auth/refresh`
* Формат запроса:
```json
{
    "refresh_token": "188ea22713b990a71f66f664534f8670151a5846f32ac04641110988d147610c"
}
```
* Формат ответа:
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTM1MDUzMDEuODQ5ODY5LCJpYXQiOjE2OTM1MDQ0MDEuODQ5ODcxMiwidXNlcm5hbWUiOiJpdmFuMDA3In0.ppetk8epzX-5PMz4WEICSiArlqEcXyhgKi8_gz-GQHo",
    "refresh_token": "d16accc31b9213361aec0eb4874aa237ba816c87a8d61b6e00eed1256f63c17f"
}
```

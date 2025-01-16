# gRPC Authorizer
## Методы
1. `CreateUser` - создание нового пользователя
* Принимает `username, password, email string`
* Создает новую запись в БД
* Пароль хранится в виде хэша
* Если пользователь с переданным `username` или `email` уже существует, то возвращает ошибку со статусом `ALREADY_EXISTS`
2. `Login` - авторизация пользователя
* Принимает `username, password string`
* Хэширует полученный пароль и сверяет с хранящимся в БД хэшем
* Если пользователь не найден или хэши не совпадают, возвращает ошибку со статусом `INVALID_ARGUMENT`
* Генерирует JWT-токен и сохраняет в БД
* Возвращает сгенерированный токен
3. `VerifyToken` - проверка токена
* Сравнивает полученный токен с хранящимся в БД
* Возвращает `true` если токены совпадают, в противном случае - `false`

## Конфигурация
Чтение конфигурации происходит из файла, переданного флагом `-c` (по умолчанию - чтение из корня проекта).

Конфигурация подключения к PostgreSQL
*  `PSQL_HOST` - default `localhost`
* `PSQL_PORT` - default `5432`
* `PSQL_DB_NAME` - default `postgres`
* `PSQL_USER` - default `postgres`
* `PSQL_PASSWORD` - default `postgres`
* `PSQL_SSL_MODE` - default `disable`
* `PSQL_CONN_TIMEOUT` - default `60` (в секундах)

Конфигурация gRPC-сервера
* `GRPC_HOST` - default `localhost`
* `GRPC_PORT` - default `9090`
* `GRPC_TIMEOUT` - default `60` (в секундах)

Конфигурация приложения
* `APP_COST_ENCODING` - default `8` (используется для хэширования паролей)
* `APP_SECRET_KEY` - default `my secret key` (используется для генерации токенов)
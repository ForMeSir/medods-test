﻿# medods-test

Cборка и запуск программы:
1. go mod init ruby
2. go build ./cmd/main.go
3. go run ./cmd/main.go
4. Поднятие базы mongodb

Как работает программа:
1.Отправляем в BODY GUID по ссылке http://localhost:3423/auth/sign-in получаем Refresh Aссess токены и их время жизни.
Пример:
{
"id":"a16fffb1-c6c6-4676-b248-7b7dde3c476d"
}
2.Отправляем в BODY Refresh токен по ссылке http://localhost:3423/auth/refresh и получаем Access токен если он валиден

# TODO APP


## Важная информация
* Этот проект создан для образовательных и развлекательных целей. Он не рассматривается как законченный и готовый к использованию продукт.
* В этом коде много плохих, странных, сомнительных, глупых итд. решений, которые сделаны НАМЕРЕННО.
* Лучшие практики и общепринятные методы построения программ сознательно нарушаются с целью попробовать другие способы и решения даже в ущерб производительности.
* SOLID/DRY/KISS/YAGNI также не соблюдаются, потому что я хотел исследовать разные возможности написания кода на го.
* Некоторые "ошибки" в базе данных (например, отсутствие внешних ключей) также были сделанны намеренно.

## Important
* This project exists for educational and recreational purposes only. It is not supposed to be a finished good-to-go product.
* This code has a lot of bad, strange, weird, dumb etc. decisions ON PURPOSE.
* Best practises and common ways to do things ignored for sake of finding another solutions even with worse perfomance.
* SOLID/DRY/KISS/YAGNI etc also broken and not respected because I wanted to explore different possibilities of programming in go.
* Some "mistakes" in database (e.g. no foreign keys) also made intentionally.


## Clone of educational API for my friend's ToDo App

Origin: https://social-network.samuraijs.com/docs

## Installation
1. Install swagger
```go
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Change host in main.go
```
// @host    your.host.name
```
3. Generate docs
```
swag init
```
Default route for swagger docs: host:port/swagger/


4. Create .env file and fill this params with your values:
```sh
# DB params
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=somename
DB_SSLMODE=disable

# HTTP Server params
HTTP_HOST=localhost
HTTP_PORT=9000

# Logger levels
APP_LOG_LEVEL=info
DB_LOG_LEVEL=silent

# Email
DOMAIN_NAME = "your frontpage domain name for verification link"
EMAIL_LOGIN = login
EMAIL_PASS = pass
EMAIL_REPLY = reply@emailservice.box
EMAIL_SERVICE = email.server.name

# Origins
ALLOWED_ORIGINS = "http://localhost, http://localhost:9090"
```

5. Build and run

## License
```
MIT
```

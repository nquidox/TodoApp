# TODO APP
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
```

5. Build and run

## License
```
MIT
```

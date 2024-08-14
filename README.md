# TODO APP
## Clone of educational API for my friend's ToDo App

Origin: https://social-network.samuraijs.com/docs

## Installation
1. Install swagger and generate docs
```sh
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

2. Create .env file and fill this params with your values:
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

3. Build and run
Default route for swagger docs: host:port/swagger/

## License
```
MIT
```

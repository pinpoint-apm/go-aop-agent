module server

go 1.16

require (
	github.com/go-redis/redis/v8 v8.11.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gorilla/mux v1.8.0
	go.mongodb.org/mongo-driver v1.5.3
	naver/app v0.0.0-00010101000000-000000000000
)

replace naver/app => ./app

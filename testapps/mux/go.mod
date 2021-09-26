module server

go 1.16

require (
	github.com/go-redis/redis/v8 v8.11.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gorilla/mux v1.8.0
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.mongodb.org/mongo-driver v1.5.3
	naver/app v0.0.0-00010101000000-000000000000
	github.com/pinpoint-apm/go-aop-agent v1.0.0
	github.com/pinpoint-apm/go-aop-agent/libs/httpClient v0.0.0-20210610105738-6027a2ff599f
	github.com/pinpoint-apm/go-aop-agent/libs/mongo v0.0.0-20210630061937-33ab55fd140c
	github.com/pinpoint-apm/go-aop-agent/libs/redisv8 v0.0.0-00010101000000-000000000000 // indirect
	github.com/pinpoint-apm/go-aop-agent/libs/sql v0.0.0-00010101000000-000000000000 // indirect
	github.com/pinpoint-apm/go-aop-agent/middleware/mux v0.0.0-20210610105738-6027a2ff599f
)

replace naver/app => ./app

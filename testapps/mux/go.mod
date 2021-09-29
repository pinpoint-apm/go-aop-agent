module server

go 1.16

require (
	github.com/go-redis/redis/v8 v8.11.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gorilla/mux v1.8.0
	github.com/pinpoint-apm/go-aop-agent v1.0.3
	github.com/pinpoint-apm/go-aop-agent/libs/httpClient v0.0.0-20210926052240-f9444c9ab5c4
	github.com/pinpoint-apm/go-aop-agent/libs/mongo v0.0.0-20210926052810-c0a16d642e12
	github.com/pinpoint-apm/go-aop-agent/libs/redisv8 v0.0.0-20210926052810-c0a16d642e12
	github.com/pinpoint-apm/go-aop-agent/libs/sql v0.0.0-20210926052810-c0a16d642e12
	github.com/pinpoint-apm/go-aop-agent/middleware/mux v0.0.0-20210926052810-c0a16d642e12
	go.mongodb.org/mongo-driver v1.5.3
	naver/app v0.0.0-00010101000000-000000000000
)

replace naver/app => ./app

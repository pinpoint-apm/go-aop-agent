module server

go 1.16

require (
	github.com/labstack/echo v3.3.10+incompatible
	github.com/pinpoint-apm/go-aop-agent v1.0.3
	github.com/pinpoint-apm/go-aop-agent/libs/httpClient v0.0.0-20210926052240-f9444c9ab5c4
	github.com/pinpoint-apm/go-aop-agent/libs/mongo v0.0.0-20210926052810-c0a16d642e12
	github.com/pinpoint-apm/go-aop-agent/libs/redisv8 v0.0.0-20211008070231-e0080eb827a0
	github.com/pinpoint-apm/go-aop-agent/libs/sql v0.0.0-20211008070231-e0080eb827a0
	github.com/pinpoint-apm/go-aop-agent/libs/transport v0.0.0-20211008070231-e0080eb827a0
	github.com/pinpoint-apm/go-aop-agent/middleware/echo v0.0.0-20210926052810-c0a16d642e12
	go.mongodb.org/mongo-driver v1.5.3
)

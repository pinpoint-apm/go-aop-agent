module server

go 1.16

require (
	github.com/labstack/echo v3.3.10+incompatible
	github.com/pinpoint-apm/go-aop-agent v1.0.0
	github.com/pinpoint-apm/go-aop-agent/libs/httpClient v0.0.0-20210610105738-6027a2ff599f
	github.com/pinpoint-apm/go-aop-agent/libs/mongo v0.0.0-20210610105738-6027a2ff599f
	github.com/pinpoint-apm/go-aop-agent/middleware/echo v0.0.0-20210610105738-6027a2ff599f
	go.mongodb.org/mongo-driver v1.5.3
)

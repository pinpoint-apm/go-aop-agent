module github.com/pinpoint-apm/go-aop-agent/middleware/echoV4

go 1.16

require (
	github.com/labstack/echo/v4 v4.3.0
	github.com/pinpoint-apm/go-aop-agent v0.0.1
)

replace github.com/pinpoint-apm/go-aop-agent => ../../

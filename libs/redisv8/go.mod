module github.com/pinpoint-apm/go-aop-agent/libs/redisv8

go 1.16

require (
	github.com/go-redis/redis/v8 v8.11.0
	github.com/pinpoint-apm/go-aop-agent v0.0.1
)

replace github.com/pinpoint-apm/go-aop-agent => ../../

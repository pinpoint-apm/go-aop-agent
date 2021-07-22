module github.com/pinpoint-apm/go-aop-agent/libs/mongo

go 1.16

require (
	go.mongodb.org/mongo-driver v1.5.3
	github.com/pinpoint-apm/go-aop-agent v0.0.1
)

replace github.com/pinpoint-apm/go-aop-agent => ../../

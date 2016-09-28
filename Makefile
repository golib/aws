all: gobuild gotest

gobuild: goclean goinstall

gorebuild: goclean goreinstall

goclean:
	go clean ./...

goinstall:
	go get github.com/golib/assert

goreinstall:
	go get -a -v github.com/golib/assert

gotest:
	go test github.com/golib/aws/service
	go test github.com/golib/aws/service/awserr
	go test github.com/golib/aws/service/awstesting
	go test github.com/golib/aws/service/awsutil
	go test github.com/golib/aws/service/client
	go test github.com/golib/aws/service/corehandlers
	go test github.com/golib/aws/service/credentials
	go test github.com/golib/aws/service/defaults
	go test github.com/golib/aws/service/endpoints
	go test github.com/golib/aws/service/request
	go test github.com/golib/aws/service/session
	go test github.com/golib/aws/service/signer/v4

travis: gobuild gotest

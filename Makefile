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
	go test github.com/golib/aws/internal
	go test github.com/golib/aws/internal/awserr
	go test github.com/golib/aws/internal/awstesting
	go test github.com/golib/aws/internal/awsutil
	go test github.com/golib/aws/internal/client
	go test github.com/golib/aws/internal/corehandlers
	go test github.com/golib/aws/internal/credentials
	go test github.com/golib/aws/internal/defaults
	go test github.com/golib/aws/internal/endpoints
	go test github.com/golib/aws/internal/request
	go test github.com/golib/aws/internal/session
	go test github.com/golib/aws/internal/signer/v4

travis: gobuild gotest

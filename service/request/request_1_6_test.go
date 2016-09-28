// +build go1.6

package request_test

import (
	"testing"

	"github.com/golib/assert"
	"github.com/golib/aws/service"
	"github.com/golib/aws/service/client"
	"github.com/golib/aws/service/client/metadata"
	"github.com/golib/aws/service/defaults"
	"github.com/golib/aws/service/endpoints"
	"github.com/golib/aws/service/request"
)

// go version 1.4 and 1.5 do not return an error. Version 1.5 will url encode
// the uri while 1.4 will not
func TestRequestInvalidEndpoint(t *testing.T) {
	endpoint, _ := endpoints.NormalizeEndpoint("localhost:80 ", "test-service", "test-region", false, false)

	r := request.New(
		service.Config{},
		metadata.ClientInfo{Endpoint: endpoint},
		defaults.Handlers(),
		client.DefaultRetryer{},
		&request.Operation{},
		nil,
		nil,
	)

	assert.Error(t, r.Error)
}

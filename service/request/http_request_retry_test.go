// +build go1.5

package request_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/golib/assert"
	"github.com/golib/aws/service"
	"github.com/golib/aws/service/awstesting/mock"
	"github.com/golib/aws/service/request"
)

func TestRequestCancelRetry(t *testing.T) {
	c := make(chan struct{})

	reqNum := 0
	s := mock.NewMockClient(service.NewConfig().WithMaxRetries(10))
	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.Clear()
	s.Handlers.UnmarshalMeta.Clear()
	s.Handlers.UnmarshalError.Clear()
	s.Handlers.Send.PushFront(func(r *request.Request) {
		reqNum++
		r.Error = errors.New("net/http: canceled")
	})
	out := &testData{}
	r := s.NewRequest(&request.Operation{Name: "Operation"}, nil, out)
	r.HTTPRequest.Cancel = c
	close(c)

	err := r.Send()
	assert.True(t, strings.Contains(err.Error(), "canceled"))
	assert.Equal(t, 1, reqNum)
}

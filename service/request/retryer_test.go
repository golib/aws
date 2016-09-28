package request

import (
	"testing"

	"github.com/golib/assert"
	"github.com/golib/aws/service/awserr"
)

func TestRequestThrottling(t *testing.T) {
	req := Request{}

	req.Error = awserr.New("Throttling", "", nil)
	assert.True(t, req.IsErrorThrottle())
}

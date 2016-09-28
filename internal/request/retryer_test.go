package request

import (
	"testing"

	"github.com/golib/assert"
	"github.com/golib/aws/internal/awserr"
)

func TestRequestThrottling(t *testing.T) {
	req := Request{}

	req.Error = awserr.New("Throttling", "", nil)
	assert.True(t, req.IsErrorThrottle())
}

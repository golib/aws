package awsutil

import (
	"testing"

	"github.com/golib/assert"
	"github.com/golib/aws/internal"
)

func TestDeepEqual(t *testing.T) {
	cases := []struct {
		a, b  interface{}
		equal bool
	}{
		{"a", "a", true},
		{"a", "b", false},
		{"a", internal.String(""), false},
		{"a", nil, false},
		{"a", internal.String("a"), true},
		{(*bool)(nil), (*bool)(nil), true},
		{(*bool)(nil), (*string)(nil), false},
		{nil, nil, true},
	}

	for i, c := range cases {
		assert.Equal(t, c.equal, DeepEqual(c.a, c.b), "%d, a:%v b:%v, %t", i, c.a, c.b, c.equal)
	}
}

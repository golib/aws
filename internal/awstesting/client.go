package awstesting

import (
	"github.com/golib/aws/internal"
	"github.com/golib/aws/internal/client"
	"github.com/golib/aws/internal/client/metadata"
	"github.com/golib/aws/internal/defaults"
)

// NewClient creates and initializes a generic service client for testing.
func NewClient(cfgs ...*internal.Config) *client.Client {
	info := metadata.ClientInfo{
		Endpoint:    "http://endpoint",
		SigningName: "",
	}
	def := defaults.Get()
	def.Config.MergeIn(cfgs...)

	return client.New(*def.Config, info, def.Handlers)
}

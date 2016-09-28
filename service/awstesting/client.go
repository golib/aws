package awstesting

import (
	"github.com/golib/aws/service"
	"github.com/golib/aws/service/client"
	"github.com/golib/aws/service/client/metadata"
	"github.com/golib/aws/service/defaults"
)

// NewClient creates and initializes a generic service client for testing.
func NewClient(cfgs ...*service.Config) *client.Client {
	info := metadata.ClientInfo{
		Endpoint:    "http://endpoint",
		SigningName: "",
	}
	def := defaults.Get()
	def.Config.MergeIn(cfgs...)

	return client.New(*def.Config, info, def.Handlers)
}

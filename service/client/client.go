package client

import (
	"fmt"

	"github.com/golib/aws/service"
	"github.com/golib/aws/service/client/metadata"
	"github.com/golib/aws/service/request"
)

// A Config provides configuration to a service client instance.
type Config struct {
	Config                  *service.Config
	Handlers                request.Handlers
	Endpoint, SigningRegion string
}

// ConfigProvider provides a generic way for a service client to receive
// the ClientConfig without circular dependencies.
type ConfigProvider interface {
	ClientConfig(serviceName string, cfgs ...*service.Config) Config
}

// A Client implements the base client request and response handling
// used by all service clients.
type Client struct {
	request.Retryer
	metadata.ClientInfo

	Config   service.Config
	Handlers request.Handlers
}

// New will return a pointer to a new initialized service client.
func New(cfg service.Config, info metadata.ClientInfo, handlers request.Handlers, options ...func(*Client)) *Client {
	svc := &Client{
		Config:     cfg,
		ClientInfo: info,
		Handlers:   handlers,
	}

	switch retryer, ok := cfg.Retryer.(request.Retryer); {
	case ok:
		svc.Retryer = retryer
	case cfg.Retryer != nil && cfg.Logger != nil:
		s := fmt.Sprintf("WARNING: %T does not implement request.Retryer; using DefaultRetryer instead", cfg.Retryer)
		cfg.Logger.Log(s)
		fallthrough
	default:
		maxRetries := service.IntValue(cfg.MaxRetries)
		if cfg.MaxRetries == nil || maxRetries == service.UseServiceDefaultRetries {
			maxRetries = 3
		}
		svc.Retryer = DefaultRetryer{NumMaxRetries: maxRetries}
	}

	svc.AddDebugHandlers()

	for _, option := range options {
		option(svc)
	}

	return svc
}

// NewRequest returns a new Request pointer for the service API
// operation and parameters.
func (c *Client) NewRequest(operation *request.Operation, params interface{}, data interface{}) *request.Request {
	return request.New(c.Config, c.ClientInfo, c.Handlers, c.Retryer, operation, params, data)
}

// AddDebugHandlers injects debug logging handlers into the service to log request
// debug information.
func (c *Client) AddDebugHandlers() {
	if !c.Config.LogLevel.AtLeast(service.LogDebug) {
		return
	}

	c.Handlers.Send.PushFront(logRequest)
	c.Handlers.Send.PushBack(logResponse)
}

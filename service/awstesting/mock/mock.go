package mock

import (
	"net/http"
	"net/http/httptest"

	"github.com/golib/aws/service"
	"github.com/golib/aws/service/client"
	"github.com/golib/aws/service/client/metadata"
	"github.com/golib/aws/service/session"
)

var (
	// server is the mock server that simply writes a 200 status back to the client
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Session is a mock session which is used to hit the mock server
	Session = session.Must(session.NewSession(&service.Config{
		Endpoint:   service.String(server.URL[7:]),
		DisableSSL: service.Bool(true),
	}))
)

// NewMockClient creates and initializes a client that will connect to the mock server
func NewMockClient(cfgs ...*service.Config) *client.Client {
	c := Session.ClientConfig("Mock", cfgs...)

	svc := client.New(
		*c.Config,
		metadata.ClientInfo{
			ServiceName:   "Mock",
			SigningRegion: c.SigningRegion,
			Endpoint:      c.Endpoint,
			APIVersion:    "2015-12-08",
			JSONVersion:   "1.1",
			TargetPrefix:  "MockServer",
		},
		c.Handlers,
	)

	return svc
}

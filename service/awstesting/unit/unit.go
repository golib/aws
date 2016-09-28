// Package unit performs initialization and validation for unit tests
package unit

import (
	"github.com/golib/aws/service"
	"github.com/golib/aws/service/credentials"
	"github.com/golib/aws/service/session"
)

// Session is a shared session for unit tests to use.
var Session = session.Must(session.NewSession(service.NewConfig().
	WithCredentials(credentials.NewStaticCredentials("AKID", "SECRET", "SESSION")).
	WithRegion("mock-region")))

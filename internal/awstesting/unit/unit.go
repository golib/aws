// Package unit performs initialization and validation for unit tests
package unit

import (
	"github.com/golib/aws/internal"
	"github.com/golib/aws/internal/credentials"
	"github.com/golib/aws/internal/session"
)

// Session is a shared session for unit tests to use.
var Session = session.Must(session.NewSession(internal.NewConfig().
	WithCredentials(credentials.NewStaticCredentials("AKID", "SECRET", "SESSION")).
	WithRegion("mock-region")))

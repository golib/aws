package v4_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/golib/assert"
	"github.com/golib/aws/service"
	"github.com/golib/aws/service/awstesting/unit"
	"github.com/golib/aws/service/client"
	"github.com/golib/aws/service/client/metadata"
	"github.com/golib/aws/service/request"
	"github.com/golib/aws/service/signer/v4"
)

var (
	testSvc = func() *client.Client {
		c := unit.Session.ClientConfig("mock")

		client := client.New(
			*c.Config,
			metadata.ClientInfo{
				ServiceName:   "mock",
				SigningRegion: c.SigningRegion,
				Endpoint:      c.Endpoint,
				APIVersion:    "2006-03-01",
			},
			c.Handlers,
		)

		client.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)

		return client
	}()
)

func TestPresignHandler(t *testing.T) {
	type (
		input struct {
			Bucket             *string
			Key                *string
			ContentDisposition *string
			ACL                *string
		}

		output struct{}
	)

	op := &request.Operation{
		Name:       "PutObject",
		HTTPMethod: "PUT",
		HTTPPath:   "/bucket/key",
	}
	in := &input{
		Bucket: service.String("bucket"),
		Key:    service.String("key"),
	}
	out := &output{}

	req := testSvc.NewRequest(op, in, out)
	req.Time = time.Unix(0, 0)
	req.HTTPRequest.Header.Add("Content-Disposition", "a+b c$d")
	req.HTTPRequest.Header.Add("X-Aws-Acl", "public-read")

	urlstr, err := req.Presign(5 * time.Minute)
	assert.NoError(t, err)

	expectedDate := "19700101T000000Z"
	expectedHeaders := "content-disposition;host;x-aws-acl"
	expectedSig := "be0a55ca135857393d3b0530b12ac10ffd6e0c3f6df4bbbc943c0faf9691f1d4"
	expectedCred := "AKID/19700101/mock-region/mock/aws4_request"

	u, _ := url.Parse(urlstr)
	urlQ := u.Query()
	assert.Equal(t, expectedSig, urlQ.Get("X-Aws-Signature"))
	assert.Equal(t, expectedCred, urlQ.Get("X-Aws-Credential"))
	assert.Equal(t, expectedHeaders, urlQ.Get("X-Aws-SignedHeaders"))
	assert.Equal(t, expectedDate, urlQ.Get("X-Aws-Date"))
	assert.Equal(t, "300", urlQ.Get("X-Aws-Expires"))

	assert.NotContains(t, urlstr, "+") // + encoded as %20
}

func TestPresignRequest(t *testing.T) {
	type (
		input struct {
			Bucket             *string
			Key                *string
			ContentDisposition *string
			ACL                *string
		}

		output struct{}
	)

	op := &request.Operation{
		Name:       "PutObject",
		HTTPMethod: "PUT",
		HTTPPath:   "/bucket/key",
	}
	in := &input{
		Bucket: service.String("bucket"),
		Key:    service.String("key"),
	}
	out := &output{}

	req := testSvc.NewRequest(op, in, out)
	req.Time = time.Unix(0, 0)
	req.HTTPRequest.Header.Add("Content-Disposition", "a+b c$d")
	req.HTTPRequest.Header.Add("X-Aws-Acl", "public-read")

	urlstr, headers, err := req.PresignRequest(5 * time.Minute)
	assert.NoError(t, err)

	expectedDate := "19700101T000000Z"
	expectedHeaders := "content-disposition;host;x-aws-acl;x-aws-content-sha256"
	expectedSig := "d2f6393819c752aa5089f994495ac996b1623288e4baa66fa587d89099b691a8"
	expectedCred := "AKID/19700101/mock-region/mock/aws4_request"
	expectedHeaderMap := http.Header{
		"x-aws-acl":            []string{"public-read"},
		"content-disposition":  []string{"a+b c$d"},
		"x-aws-content-sha256": []string{"UNSIGNED-PAYLOAD"},
	}

	u, _ := url.Parse(urlstr)
	urlQ := u.Query()
	assert.Equal(t, expectedSig, urlQ.Get("X-Aws-Signature"))
	assert.Equal(t, expectedCred, urlQ.Get("X-Aws-Credential"))
	assert.Equal(t, expectedHeaders, urlQ.Get("X-Aws-SignedHeaders"))
	assert.Equal(t, expectedDate, urlQ.Get("X-Aws-Date"))
	assert.Equal(t, expectedHeaderMap, headers)
	assert.Equal(t, "300", urlQ.Get("X-Aws-Expires"))

	assert.NotContains(t, urlstr, "+") // + encoded as %20
}

package client

import (
	"fmt"
	"io/ioutil"
	"net/http/httputil"

	"github.com/golib/aws/internal"
	"github.com/golib/aws/internal/request"
)

const (
	logReqMsg = `DEBUG: Request %s/%s Details:
---[ REQUEST POST-SIGN ]-----------------------------
%s
-----------------------------------------------------`

	logReqErrMsg = `DEBUG ERROR: Request %s/%s:
---[ REQUEST DUMP ERROR ]-----------------------------
%s
-----------------------------------------------------`
)

func logRequest(r *request.Request) {
	logBody := r.Config.LogLevel.Matches(internal.LogDebugWithHTTPBody)
	dumpedBody, err := httputil.DumpRequestOut(r.HTTPRequest, logBody)
	if err != nil {
		r.Config.Logger.Log(fmt.Sprintf(logReqErrMsg, r.ClientInfo.ServiceName, r.Operation.Name, err))
		return
	}

	if logBody {
		// Reset the request body because dumpRequest will re-wrap the r.HTTPRequest's
		// Body as a NoOpCloser and will not be reset after read by the HTTP
		// client reader.
		r.Body.Seek(r.BodyStart, 0)
		r.HTTPRequest.Body = ioutil.NopCloser(r.Body)
	}

	r.Config.Logger.Log(fmt.Sprintf(logReqMsg, r.ClientInfo.ServiceName, r.Operation.Name, string(dumpedBody)))
}

const (
	logRespMsg = `DEBUG: Response %s/%s Details:
---[ RESPONSE ]--------------------------------------
%s
-----------------------------------------------------`

	logRespErrMsg = `DEBUG ERROR: Response %s/%s:
---[ RESPONSE DUMP ERROR ]-----------------------------
%s
-----------------------------------------------------`
)

func logResponse(r *request.Request) {
	var msg = "no response data"
	if r.HTTPResponse != nil {
		logBody := r.Config.LogLevel.Matches(internal.LogDebugWithHTTPBody)
		dumpedBody, err := httputil.DumpResponse(r.HTTPResponse, logBody)
		if err != nil {
			r.Config.Logger.Log(fmt.Sprintf(logRespErrMsg, r.ClientInfo.ServiceName, r.Operation.Name, err))
			return
		}

		msg = string(dumpedBody)
	} else if r.Error != nil {
		msg = r.Error.Error()
	}
	r.Config.Logger.Log(fmt.Sprintf(logRespMsg, r.ClientInfo.ServiceName, r.Operation.Name, msg))
}

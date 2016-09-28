package awsutil

import (
	"net/url"
	"strings"
)

// EscapePath escapes part of a URL path in Amazon style
func EscapePath(path string, encodeSep bool) string {
	pieces := strings.Split(path, "/")
	if len(pieces) == 1 {
		s := strings.Replace(url.QueryEscape(path), "+", "%20", -1)
		if encodeSep {
			s = strings.Replace(s, "/", "%2F", -1)
		}

		return s
	}

	for i, seg := range pieces {
		pieces[i] = url.QueryEscape(seg)
	}

	sep := "/"
	if encodeSep {
		sep = "%2F"
	}
	return strings.Replace(strings.Join(pieces, sep), "+", "%20", -1)
}

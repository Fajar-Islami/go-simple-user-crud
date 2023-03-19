package rest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/container"
)

// Transport custom transport to use a logger, implement http.RoundTripper
type Transport struct {
	Transport   http.RoundTripper
	Log         *container.Logger
	IsDebugging bool
}

// TransportWithLogger wrap transport for logging
func TransportWithLogger(opt Opts) http.RoundTripper {
	tr := http.DefaultTransport

	return &Transport{
		Transport:   tr,
		Log:         &opt.Logger,
		IsDebugging: opt.IsDebugging,
	}
}

// RoundTrip custom roundtrip implements http.RoundTripper
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	transport := http.DefaultTransport
	if t.Transport != nil {
		transport = t.Transport
	}

	reqBody, _ := hookRequest(req)

	resp, err := transport.RoundTrip(req)
	duration := time.Since(start)
	if err != nil {
		return resp, err
	}

	resBody, _ := hookResponse(resp)

	if t.Log != nil {
		t.Log.Log.Info().
			Str("Request.URI", req.Method+" "+req.URL.String()).
			Dur("Request.Latency", duration).
			Str("Request.header", stringifyHeader(req.Header)).
			Str("Request.body", string(reqBody)).
			Str("Response.header", stringifyHeader(resp.Header)).
			Str("Response.body", string(resBody)).
			Msg("HTTP Log")

	}

	if t.IsDebugging {
		fmt.Printf(`
----------------------------------------------------------------
HTTP Request 
Date: %s
Latency: %f seconds
----------------------------------------------------------------
%s
%s
%s
%s
----------------------------------------------------------------
HTTP Response
----------------------------------------------------------------
%s
%s

`, start.Format("Mon, 2 Jan 2006 15:04:05 MST"),
			duration.Seconds(),
			req.Method,
			req.URL.String(),
			curlifyHeader(req.Header),
			curlifyBody(reqBody),
			curlifyHeader(resp.Header),
			curlifyBody(resBody))
	}

	return resp, err
}

func hookRequest(c *http.Request) (body []byte, err error) {
	if c.Body != nil { // Read
		body, err = ioutil.ReadAll(c.Body)
		if err != nil {
			return body, err
		}
	}
	c.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, err

}

func hookResponse(c *http.Response) (body []byte, err error) {
	if c.Body != nil { // Read
		body, err = ioutil.ReadAll(c.Body)
		if err != nil {
			return body, err
		}
	}
	c.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, err
}

func stringifyHeader(headers http.Header) (h string) {
	if headers != nil {
		var temp []string
		for k, v := range headers {
			temp = append(temp, fmt.Sprintf("%s: %s", k, strings.Join(v, " ")))
		}
		return strings.Join(temp, ", ")
	}
	return
}

func curlifyHeader(headers http.Header) (h string) {
	if headers != nil {
		var temp []string
		for k, v := range headers {
			temp = append(temp, fmt.Sprintf(`-H '%s: %s'`, k, strings.Join(v, " ")))
		}
		return strings.Join(temp, "\n")
	}
	return
}

func curlifyBody(body []byte) string {
	if body == nil {
		return ""
	}
	return fmt.Sprintf(`-d %s`, string(body))
}

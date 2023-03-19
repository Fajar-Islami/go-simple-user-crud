package rest

import (
	"io"
	"net/http"
	"time"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/container"
)

// Opts rest options
type Opts struct {
	Timeout     time.Duration
	Logger      container.Logger
	IsDebugging bool
}

func DoRequest(method, uri string, opts Opts, body ...io.Reader) (*http.Response, error) {
	var err error
	var client = &http.Client{
		Transport: TransportWithLogger(opts),
		Timeout:   opts.Timeout,
	}
	var newBody io.Reader

	if len(body) < 1 {
		newBody = nil
	} else {
		newBody = body[0]
	}

	request, err := http.NewRequest(method, uri, newBody)
	if err != nil {
		return nil, err
	}

	return client.Do(request)
}

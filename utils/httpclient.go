package utils

import (
	"net/http"
	"time"
)

func NewHttpTransport() *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxIdleConnsPerHost = 100
	return t
}

func NewHttpClient() *http.Client{
	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: NewHttpTransport(),
	}
	return client
}


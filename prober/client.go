package prober

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	co sync.Once
	c  *http.Client
)

func httpClient() *http.Client {
	co.Do(func() {
		c = &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout:   5 * time.Second,
				ResponseHeaderTimeout: 5 * time.Second,
				ExpectContinueTimeout: 2 * time.Second,
			},
		}
	})
	return c
}

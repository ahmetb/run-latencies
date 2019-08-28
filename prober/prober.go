package prober

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"
)

type ErrorKind string

const (
	RequestError     ErrorKind = "request"
	MeasurementError ErrorKind = "measurement"
	InputError       ErrorKind = "input"
)

type Stage struct {
	Name       string `json:"name"`
	DurationMs int    `json:"duration_ms"`
}

type ErrorOutcome struct {
	Kind    ErrorKind `json:"kind"`
	Message string    `json:"message"`
}

type Measurement struct {
	Stages     []Stage     `json:"stages,omitempty"`
	StatusCode int         `json:"status_code,omitempty"`
	TotalMs    int         `json:"total_ms,omitempty"`
}

func Probe(u *url.URL) (*Measurement, *ErrorOutcome) {
	hc := httpClient()
	start := time.Now()

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("connection", "close")
	var (
		dnsStart    time.Time
		dnsDuration time.Duration

		connectStart    time.Time
		connectDuration time.Duration

		tlsHandshakeStart    time.Time
		tlsHandshakeDuration time.Duration

		wroteHeaders      time.Time
		firstByteDuration time.Duration

		downloadStart    time.Time
		downloadDuration time.Duration
	)

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			log.Printf("dns start")
			dnsStart = time.Now()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			log.Printf("dns done")
			dnsDuration = time.Since(dnsStart)
		},
		ConnectStart: func(_, _ string) {
			log.Printf("connect start")
			connectStart = time.Now()
		},
		ConnectDone: func(_, _ string, _ error) {
			log.Printf("connect done")
			connectDuration = time.Since(connectStart)
		},
		TLSHandshakeStart: func() {
			log.Printf("tls start")
			tlsHandshakeStart = time.Now()
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			log.Printf("tls done")
			tlsHandshakeDuration = time.Since(tlsHandshakeStart)
		},
		WroteHeaders: func() {
			log.Printf("wrote headers")
			wroteHeaders = time.Now()
		},
		GotFirstResponseByte: func() {
			log.Printf("first response byte")
			firstByteDuration = time.Since(wroteHeaders)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := hc.Do(req)
	if err != nil {
		return nil, &ErrorOutcome{RequestError, err.Error()}
	}
	defer resp.Body.Close()

	downloadStart = time.Now()
	ioutil.ReadAll(resp.Body)
	downloadDuration = time.Since(downloadStart)

	total := time.Since(start)
	log.Printf("measured %s, total=%v", u.String(), total)

	if dnsDuration == 0 {
		return nil, &ErrorOutcome{MeasurementError, "no dns resolution was done"}
	}
	if connectDuration == 0 {
		return nil, &ErrorOutcome{MeasurementError, "no tcp connection was done"}
	}
	if tlsHandshakeDuration == 0 {
		return nil, &ErrorOutcome{MeasurementError, "no tls handshake was done"}
	}

	return &Measurement{
		StatusCode: resp.StatusCode,
		TotalMs:    int(total / time.Millisecond),
		Stages: []Stage{
			{Name: "dns",
				DurationMs: int(dnsDuration / time.Millisecond)},
			{Name: "tcp",
				DurationMs: int(connectDuration / time.Millisecond)},
			{Name: "tls",
				DurationMs: int(tlsHandshakeDuration / time.Millisecond)},
			{Name: "first_byte",
				DurationMs: int(firstByteDuration / time.Millisecond)},
			{Name: "download",
				DurationMs: int(downloadDuration / time.Millisecond)},
		},
	}, nil
}

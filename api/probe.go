package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runprobe/prober"
	"strings"
)


type Resp struct {
	Success bool                 `json:"success"`
	Data    *prober.Measurement  `json:"data,omitempty"`
	Error   *prober.ErrorOutcome `json:"error,omitempty"`
}

func init() { log.SetFlags(log.Lmicroseconds) }

func Probe(w http.ResponseWriter, req *http.Request) {
	uv := req.URL.Query().Get("url")
	if uv == "" {
		respond(w, Resp{
			Error: &prober.ErrorOutcome{
				Kind:    prober.InputError,
				Message: "url is not specified"}})
		return
	}
	u, err := url.Parse(uv)
	if err != nil {
		respond(w, Resp{
			Error: &prober.ErrorOutcome{
				Kind:    prober.InputError,
				Message: "failed to parse url parameter"}})
		return
	}
	if !safeURL(u) {
		respond(w, Resp{
			Error: &prober.ErrorOutcome{
				Kind:    prober.InputError,
				Message: "url hostname is not whitelisted"}})
		return
	}

	v, e := prober.Probe(u)
	respond(w, Resp{
		Success: err == nil,
		Data:    v,
		Error:   e,
	})
}

func safeURL(u *url.URL) bool {
	if u.Hostname() == "localhost" {
		return true
	}
	if u.Scheme != "https" {
		return false
	}
	return strings.HasSuffix(u.Hostname(), ".run.app")
}

func respond(w http.ResponseWriter, r Resp) {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(r)
}

func Home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok")
}

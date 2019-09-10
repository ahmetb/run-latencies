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

func init() { log.SetFlags(log.Lmicroseconds) }

var _ http.HandlerFunc = Probe // type checks

func Probe(w http.ResponseWriter, req *http.Request) {
	uv := req.URL.Query().Get("url")
	if uv == "" {
		respond(w, prober.ProbeResp{
			Error: &prober.ErrorOutcome{
				Kind:    prober.InputError,
				Message: "url is not specified"}})
		return
	}
	u, err := url.Parse(uv)
	if err != nil {
		respond(w, prober.ProbeResp{
			Error: &prober.ErrorOutcome{
				Kind:    prober.InputError,
				Message: "failed to parse url parameter"}})
		return
	}
	if !safeURL(u) {
		respond(w, prober.ProbeResp{
			Error: &prober.ErrorOutcome{
				Kind:    prober.InputError,
				Message: "url hostname is not whitelisted"}})
		return
	}

	w.Header().Set("netgo", fmt.Sprintf("%v", prober.Netgo))
	v, e := prober.Probe(u)
	respond(w, prober.ProbeResp{
		Netgo:   prober.Netgo,
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

func respond(w http.ResponseWriter, r prober.ProbeResp) {
	w.Header().Set("content-type", "application/json")
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(r)
}

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"runprobe/prober"
)

var (
	flURL string
)

func init() {
	flag.StringVar(&flURL, "url", "", "url to probe")
	flag.Parse()
}

func main() {
	if flURL == "" {
		log.Fatal("-url not specified")
	}
	u, err := url.Parse(flURL)
	if err != nil {
		log.Fatal(err)
	}
	v, e := prober.Probe(u)
	if e != nil {
		log.Fatalf("error: kind=%v msg=%s", e.Kind, e.Message)
	}
	n := json.NewEncoder(os.Stdout)
	n.SetIndent("", "  ")
	n.Encode(v)
}

package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	mu     sync.Mutex
	served bool

	respBody = bytes.Repeat([]byte{'a', 'b', 'c'}, 1024)
)

func main() {
	addr := os.Getenv("ADDR")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env var not set")
	}
	http.HandleFunc("/", home)
	http.ListenAndServe(addr+":"+port, nil)
}

func home(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %s not allowed", req.Method)
		return
	}
	mu.Lock()
	warm := served
	served = true
	mu.Unlock()

	if warm {
		w.Header().Set("request-state", "warm")
	} else {
		w.Header().Set("request-state", "cold")
	}

	if req.Method == http.MethodGet {
		w.Write(respBody)
	}
}

package run_latencies

import (
	"log"
	"net/http"
	"os"
	handler "runprobe/api"
)

func main() {
	addr := os.Getenv("ADDR")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env var not set")
	}
	listenAddr := addr + ":" + port
	http.HandleFunc("/probe", handler.Probe)
	log.Printf("server starting on " + listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

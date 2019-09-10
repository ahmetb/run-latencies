// +build netgo

package prober

import "net/http"

func init() {
	Netgo = true
}

func Dummy(w http.ResponseWriter, req *http.Request) {}

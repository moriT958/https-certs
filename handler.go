package main

import (
	"fmt"
	"net/http"
)

type customHandler struct{}

func (h *customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from http2!")
}

package http

import (
	"fmt"
	"net/http"
)

type HTTPHandler struct {
	// Add any necessary fields here
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle the HTTP request
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}

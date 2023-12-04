package http

import (
	"net/http"
)

type HTTPServer struct {
	Port     int16
	Handlers []*http.Handler
}

package api

import (
	"fmt"
	"net/http"
	"time"

	myLog "github.com/subbbbbaru/first-sample/pkg/log"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port int, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + fmt.Sprintf("%d", port),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	myLog.Info().Printf("Server run at %s port\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

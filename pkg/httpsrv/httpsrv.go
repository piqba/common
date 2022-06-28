package httpsrv

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/piqba/common/pkg/logger"
)

type server struct {
	*http.Server
	log *log.AppLogger
}

// NewServer return a custom server without SSL
func NewServer(host string, port int, mux http.Handler, logger *log.AppLogger, write, read, idle time.Duration) *server {
	newLogger := logger
	if logger == nil {
		newLogger = log.NewLogger("INFO")
	}
	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: mux,
		// Good practice: enforce timeouts for servers you create!
		ReadTimeout:  read,
		WriteTimeout: write,
		IdleTimeout:  idle,
	}
	return &server{s, newLogger}
}

// NewServerSSL return a custom server with SSL
func NewServerSSL(host string, port int, mux http.Handler, logger *log.AppLogger, write, read, idle time.Duration) *server {
	newLogger := logger
	if logger == nil {
		newLogger = log.NewLogger("INFO")
	}
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	s := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf("%s:%d", host, port),
		// Good practice: enforce timeouts for servers you create!
		ReadTimeout:  read,
		WriteTimeout: write,
		IdleTimeout:  idle,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	return &server{s, newLogger}
}

// Start runs ListenAndServe on the http.Server with graceful shutdown
func (srv *server) Start() {
	srv.log.Info(
		" starting server...",
		//log.F{"ID": "hola"}, // this and example
	)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			srv.log.Error(errors.Errorf("could not listen on %s due to %s", srv.Addr, err).Error())
		}
	}()
	srv.log.Info(fmt.Sprintf(" server is ready to handle requests %s", srv.Addr))
	srv.gracefulShutdown()
}

// StartSSL runs ListenAndServe on the http.Server with graceful shutdown
func (srv *server) StartSSL(crt, key string) {
	srv.log.Info(" starting server...")

	go func() {
		if err := srv.ListenAndServeTLS(crt, key); err != nil && err != http.ErrServerClosed {
			srv.log.Error(errors.Errorf("could not listen on %s due to %s", srv.Addr, err).Error())
		}
	}()
	srv.log.Info(fmt.Sprintf(" server is ready to handle requests %s", srv.Addr))
	srv.gracefulShutdown()
}

// gracefulShutdown ...
func (srv *server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	srv.log.Info(fmt.Sprintf(" server is shutting down %s", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		srv.log.Error(errors.Errorf("could not gracefully shutdown the server %s", err).Error())

	}
	srv.log.Info(" server stopped")
}

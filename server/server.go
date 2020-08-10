package server

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pu "github.com/nudelfabrik/portUpdate"
	"github.com/nudelfabrik/portUpdate/server/templates"
)

type Server struct {
	entrys    []pu.Entry
	templates *template.Template
	enableSSL bool
}

type Template struct {
	Entrys    []pu.Entry
	ShowError string
}

const (
	HTMLShow = ""
	HTMLHide = "hidden"
)

func NewServer(entrys []pu.Entry) (*Server, error) {
	srv := Server{entrys: entrys}
	channel := make(chan *template.Template)
	srv.templates = templates.WatchTemplates(channel)

	go func() {
		for {
			srv.templates = <-channel
		}
	}()
	log.Println("Watching Templates in debug mode")

	return &srv, nil
}

func (srv *Server) Start() {
	// Create custom ServeMux
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/coffee", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "I'm a Teapot", http.StatusTeapot)
	})
	serveMux.HandleFunc("/", srv.IndexHandler)

	var (
		rTimeout time.Duration = 5 * time.Second

		wTimeout time.Duration = 10 * time.Second

		iTimeout time.Duration = 120 * time.Second
	)

	var tlsConf *tls.Config = nil

	if srv.enableSSL {
		tlsConf = tlsConfig("certfile", "keyfile")
	}

	httpServer := &http.Server{
		Addr:      ":8000",
		Handler:   serveMux,
		TLSConfig: tlsConf,
		// Added Timeouts to prevent resource exhaustion
		ReadTimeout:  rTimeout,
		WriteTimeout: wTimeout,
		IdleTimeout:  iTimeout,
	}

	// Setup Shutdown Signal Handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("Shutting down server...")

		if err := httpServer.Close(); err != nil {
			log.Fatalf("could not shutdown: %v", err)
		}

		os.Stdout.WriteString("\n")
	}()

	log.Println("Start Server")
	// Certs are loaded into tlsConf
	// httpServer.ListenAndServeTLS("", "")
	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func tlsConfig(certfile, keyfile string) *tls.Config {
	// Load Certificate
	cer, err := tls.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		log.Printf("Error Loading Certificates: %v\n", err)

		return nil
	}

	// TLS Config
	// https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/
	tlsConf := &tls.Config{
		Certificates:             []tls.Certificate{cer},
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519, // Go 1.8 only
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	return tlsConf
}

/*
func hstsHandler(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; preload")
		w.Header().Set("x-frame-options", "SAMEORIGIN")
		fn(w, r)
	})
}
*/

func (srv *Server) IndexHandler(w http.ResponseWriter, req *http.Request) {
	t := Template{Entrys: srv.entrys[:10]}

	err := srv.templates.ExecuteTemplate(w, "list.html", t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

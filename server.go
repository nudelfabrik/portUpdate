package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start() {

	// Create custom ServeMux
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/coffee", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "I'm a Teapot", http.StatusTeapot)
	})

	// Load Certificate
	/*cer, err := tls.LoadX509KeyPair(, srv.cfg.Keyfile())
	if err != nil {
		log.Println("Error Loading Certificates: %v", err)
		return
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
	*/

	httpServer := &http.Server{
		Addr:    ":8000",
		Handler: serveMux,
		//TLSConfig: tlsConf,
		// Added Timeouts to prevent resource exhaustion
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
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
	//httpServer.ListenAndServeTLS("", "")
	httpServer.ListenAndServe()
}

func hstsHandler(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; preload")
		w.Header().Set("x-frame-options", "SAMEORIGIN")
		fn(w, r)
	})
}

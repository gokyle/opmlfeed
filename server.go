package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	OPMLFEED_SSL_KEY = os.Getenv("OPMLFEED_SSLKEY")
	OPMLFEED_SSL_CERT = os.Getenv("OPMLFEED_SSLCERT")
	SERVER_ADDR = os.Getenv("SERVER_ADDR")
	SERVER_PORT = os.Getenv("SERVER_PORT")
	InitMux()
}

func main() {
	cfg := new(tls.Config)
	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", SERVER_ADDR, SERVER_PORT),
		Handler:        opml_mux,
		ReadTimeout:    3 * time.Second,
		WriteTimeout:   3 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      cfg,
	}
	log.Println("[+] server goes online")
	log.Fatal(srv.ListenAndServeTLS(OPMLFEED_SSL_CERT, OPMLFEED_SSL_KEY))
}

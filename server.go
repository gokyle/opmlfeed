package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const OPMLFEED_VERSION = "0.9.2"

var (
	OPMLFEED_SSL_KEY  string
	OPMLFEED_SSL_CERT string
	SERVER_ADDR       string
	SERVER_PORT       string
)

func init() {
	OPMLFEED_SSL_KEY = os.Getenv("OPMLFEED_SSLKEY")
	OPMLFEED_SSL_CERT = os.Getenv("OPMLFEED_SSLCERT")
	SERVER_ADDR = os.Getenv("SERVER_ADDR")
	SERVER_PORT = os.Getenv("SERVER_PORT")
	initMux()
}

func main() {
	crt, err := tls.LoadX509KeyPair(OPMLFEED_SSL_KEY, OPMLFEED_SSL_CERT)
	if err != nil {
		log.Println("[!] error loading SSL keypair: ", err.Error())
	}
	cfg := &tls.Config{
		Certificates: []tls.Certificate{crt},
	}
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

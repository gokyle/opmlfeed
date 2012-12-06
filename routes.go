package main

import (
        //"github.com/gokyle/webshell"
        "net/http"
)

func topRouter(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
                getRouter(w, r)
                return
        } else if r.Method == "POST" {
                postRouter(w, r)
                return
        }
}

func getRouter(w http.ResponseWriter, r *http.Request) {
        path := r.URL.Path[1:]
        if path == "" || path == "index.html" {
                serveIndex(w, r)
                return
        }

        serveErrorPage(http.StatusNotFound, w, r)
}

func postRouter(w http.ResponseWriter, r *http.Request) {
        serveErrorPage(http.StatusNotImplemented, w, r)
}

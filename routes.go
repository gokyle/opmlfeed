package main

import (
        "net/http"
)

func topRouter(w http.ResponseWriter, r *http.Request) {
        if r.URL.RawQuery=="logout" {
                store.DestroySession(r)
        }
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

        if err := servePage(w, r); err != nil {
                serveErrorPage(http.StatusNotFound, w, r)
        }
}

func postRouter(w http.ResponseWriter, r *http.Request) {
        serveErrorPage(http.StatusNotImplemented, w, r)
}

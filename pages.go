package main

import (
        "github.com/gokyle/webshell"
        "net/http"
)

type ErrorPage struct {
        Status int
        StatusText string
        StatusPicture bool
        StatusPictureFile string
}

var (
        error_tpl = webshell.MustCompileTemplate("templates/error.html")
        header_tpl = webshell.MustCompileTemplate("templates/header.html")
        footer_tpl = webshell.MustCompileTemplate("templates/footer.html")
        index_tpl = webshell.MustCompileTemplate("templates/index.html")
        header_master []byte
        footer_master []byte
)

func init() {
        var err error
        header_master, err = webshell.ServeTemplate(header_tpl,
                struct { Title string }{appname})
        if err != nil {
                panic("couldn't build templated header: " + err.Error())
        }
        footer_master, err = webshell.ServeTemplate(footer_tpl, nil)
        if err != nil {
                panic("couldn't build templated footer: " + err.Error())
        }
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
        page := struct {
                      Title string
        }{appname}
        index, err := webshell.ServeTemplate(index_tpl, page)
        if err != nil { }

        index_page := header_master
        index_page = append(index_page, index...)
        index_page = append(index_page, footer_master...)
        w.Write(index_page)
}

func serveErrorPage(status int, w http.ResponseWriter, r *http.Request) {
        errpage := ErrorPage{status, http.StatusText(status), false, ""}
        if status == http.StatusNotFound {
                errpage.StatusPicture = true
                errpage.StatusPictureFile = "assets/img/404.png"
        }
        body, err := webshell.ServeTemplate(error_tpl, errpage)
        if err != nil {
                webshell.Error500(err.Error(), "text/plain", w, r)
                return
        }
        page := header_master
        page = append(page, body...)
        page = append(page, footer_master...)
        w.Header().Set("content-type", "text/html")
        w.WriteHeader(status)
        w.Write(page)
}

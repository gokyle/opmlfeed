package main

import (
	"github.com/gokyle/webshell"
	"net/http"
)

type ErrorPage struct {
	Status            int
	StatusText        string
	StatusHaiku       string
	StatusPicture     bool
	StatusPictureFile string
}

var statusHaiku map[int]string

var (
	error_tpl     = webshell.MustCompileTemplate("templates/error.html")
	header_tpl    = webshell.MustCompileTemplate("templates/header.html")
	footer_tpl    = webshell.MustCompileTemplate("templates/footer.html")
	index_tpl     = webshell.MustCompileTemplate("templates/index.html")
	header_master []byte
	footer_master []byte
)

func init() {
	var err error
	header_master, err = webshell.BuildTemplate(header_tpl,
		struct{ Title string }{appname})
	if err != nil {
		panic("couldn't build templated header: " + err.Error())
	}
	footer_master, err = webshell.BuildTemplate(footer_tpl, nil)
	if err != nil {
		panic("couldn't build templated footer: " + err.Error())
	}

	statusHaiku = make(map[int]string, 0)
	statusHaiku[http.StatusNotFound] = `Page could not be found. Perhaps you made a typo? Is the server down?`
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	page := struct {
		Title string
	}{appname}
	index, err := webshell.BuildTemplate(index_tpl, page)
	if err != nil {
	}

	index_page := header_master
	index_page = append(index_page, index...)
	index_page = append(index_page, footer_master...)
	w.Write(index_page)
}

func servePage(w http.ResponseWriter, r *http.Request) error {
	var page = struct {
		Title string
	}{appname}
	body, err := webshell.BuildTemplateFile(`templates` + r.URL.Path, page)
	if err != nil {
		return err
	}
	full_page := header_master
	full_page = append(full_page, body...)
	full_page = append(full_page, footer_master...)
	w.Write(full_page)
        return nil
}

func serveErrorPage(status int, w http.ResponseWriter, r *http.Request) {
	errpage := ErrorPage{status, http.StatusText(status), "", false, ""}
        if haiku, ok := statusHaiku[status]; ok {
                errpage.StatusHaiku = haiku
        }
	if status == http.StatusNotFound {
		errpage.StatusPicture = true
		errpage.StatusPictureFile = "assets/img/404.png"
	}
	body, err := webshell.BuildTemplate(error_tpl, errpage)
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

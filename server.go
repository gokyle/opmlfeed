package main

import (
	"github.com/gokyle/goconfig"
	"github.com/gokyle/webshell"
	"github.com/gokyle/webshell/auth"
        "github.com/gokyle/webshell/assetcache"
        "log"
)

const config_file = "conf/opmldev.conf"
const appname = "opmlfeed2"

type AppConf struct {
	Host    string
	Port    string
	TlsKey  string
	TlsCert string
	Secure  bool
}

var (
	app   *webshell.WebApp
	store *auth.SessionStore
)

func init() {
	// auth.LookupCredentials = LookupUser
}

func main() {
	appconf := AppConf{}
	appconf.Port = "8080"

	conf, err := config.ParseFile(config_file)
        if err != nil {
                log.Println("error parsing config file: ", err.Error())
        }
	configure(conf, &appconf)

	if appconf.Secure {
		app = webshell.NewTLSApp(appname, appconf.Host, appconf.Port,
			appconf.TlsKey, appconf.TlsCert)
	} else {
		app = webshell.NewApp(appname, appconf.Host, appconf.Port)
	}
        err = assetcache.BackgroundAttachAssetCache(app, "/assets/", "assets/")
        if err != nil {
                log.Println("could not set up asset cache: ", err.Error())
                app.StaticRoute("/assets/", "assets/")
        }
        app.AddRoute("/", topRouter)
	app.Serve()
}

func configure(conf config.ConfigMap, appconf *AppConf) {
	if _, ok := conf["server"]; !ok {
		return
	}

	server := conf["server"]
	if host := server["host"]; host != "" {
		appconf.Host = host
	}

	if port := server["port"]; port != "" {
		appconf.Port = port
	}

	if key := server["tls_key"]; key != "" {
		appconf.TlsKey = key
	}

	if cert := server["tls_cert"]; cert != "" {
		appconf.TlsCert = cert
	}
}

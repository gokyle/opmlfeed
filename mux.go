package main

import (
	"bitbucket.org/kisom/opmlfeed/shorten"
	"encoding/json"
	"fmt"
	"github.com/simonz05/godis/redis"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const index_file = "index.txt"
var opml_mux = http.NewServeMux()
var index_text []byte

// regular expressions
var (
	uuid_regex_str = "\\w{8}-\\w{4}-\\w{4}-\\w{4}-\\w{12}"
	uuid_regex     = regexp.MustCompile("^" + uuid_regex_str + "$")
	uuid_strip     = regexp.MustCompile("^/f/(" + uuid_regex_str + ")$")
	feed_strip     = regexp.MustCompile("^/(\\w{6})$")
	url_regex_str  = "^(https?:\\/\\/)?([\\da-z\\.-]+)\\." +
		"([a-z\\.]{2,6})([\\/\\w \\.-]*)*\\/?$"
	url_regex = regexp.MustCompile(url_regex_str)
)

var (
	OPML_REDIS_ADDR = "127.0.0.1:6379"
	OPML_REDIS_DB   = 0
)

// type Feed represents a single feed
type Feed struct {
	Title string `json="title"`
	Xml   string `json="xml"`
	Web   string `json="web"`
}

// type Update represents a client's data as uploaded to the server
type Update struct {
	UUID  string `json="uuid"`
	Feeds []Feed `json="feeds"`
}

type Response struct {
        Feeds []Feed `json="feeds"`
}

func loadUpdate(jsonData []byte) (update *Update, err error) {
	update = new(Update)
	err = json.Unmarshal(jsonData, &update)
	return
}

func initMux() {
        var err error
        index_text, err = ioutil.ReadFile(index_file)
        if err != nil {
                log.Panic("[!] unable to load index file: ", err.Error())
        }
	opml_mux.HandleFunc("/", opmlRoot)
}

func opmlRoot(w http.ResponseWriter, r *http.Request) {
	addHStatus := func(code int) string {
		return fmt.Sprintf(" %d", code)
	}
	status := fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method,
		r.URL.String())
	if r.Method == "POST" {
		var rawJson []byte
		var err error
		if rawJson, err = ioutil.ReadAll(r.Body); err != nil {
			res := http.StatusBadRequest
			status += addHStatus(res)
			w.WriteHeader(res)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte("400 Bad Request\n"))
			w.Write([]byte(err.Error()))
			log.Println(status)
			return
		}
		update, err := loadUpdate(rawJson)
		if err != nil || !update.Validate() {
			res := http.StatusBadRequest
			status += fmt.Sprintf(" %d", res)
			status += addHStatus(res)
			w.WriteHeader(res)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte("400 Internal Server Error\n"))
			if err != nil {
				w.Write([]byte(err.Error()))
			} else {
				w.Write([]byte("invalid client data"))
			}
			log.Println(status)
			return
		}
		shortid, err := clientUpdate(update)
		if err != nil {
			res := http.StatusInternalServerError
			status += addHStatus(res)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte("500 Internal Server Error\n"))
			w.Write([]byte(err.Error()))
			log.Println(status)
		} else {
			res := http.StatusOK
			status += addHStatus(res)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte(shortid))
			log.Println(status)
		}
		return
	} else if r.Method == "GET" {
		shortid := feed_strip.ReplaceAllString(r.URL.Path, "$1")
		if len(shortid) == shorten.ShortLen {
			var jsonResp []byte
			uuid, err := shortIdToUUID(shortid)
			if len(uuid) == 0 {
				res := http.StatusNotFound
				status += addHStatus(res)
				w.WriteHeader(res)
				w.Header().Set("content-type",
					"text/plain")
				w.Write([]byte("feed not found\n"))
				log.Println(status)
				return
			}
			update, err := fetchFeed(uuid)
			if err != nil {
				res := http.StatusNotFound
				status += addHStatus(res)
				w.WriteHeader(res)
				w.Header().Set("content-type", "text/plain")
				w.Write([]byte("feed not found\n"))
				w.Write([]byte(err.Error()))
				log.Println(status)
				return
			}
                        var response Response
                        response.Feeds = update.Feeds
			jsonResp, err = json.Marshal(response)
			if err != nil {
				res := http.StatusBadRequest
				status += addHStatus(res)
				w.WriteHeader(res)
				w.Header().Set("content-type", "text/plain")
				w.Write([]byte("invalid feed list"))
				log.Println(status)
				return
			} else {
				status += addHStatus(http.StatusOK)
				w.Write(jsonResp)
			}
			log.Println(status)
			return
		} else if r.URL.Path == "/" {
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte(fmt.Sprintf("opmlfeed server %s\n\n",
				OPMLFEED_VERSION)))
                        w.Write(index_text)
			status += addHStatus(http.StatusOK)
                        log.Println(status)
		} else {
			res := http.StatusNotFound
			status += addHStatus(res)
			w.WriteHeader(res)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte("not found\n"))
			log.Println(status)
		}
	}
}

func fetchFeed(uuid string) (update *Update, err error) {
	var opml []byte
	r := redis.New("", OPML_REDIS_DB, "")

	opml, err = r.Get("OF_" + uuid)
	if len(opml) == 0 {
		return
	}
	if err == nil {
		err = json.Unmarshal(opml, &update)
	}
	return
}

func (update *Update) Validate() (valid bool) {
	valid = uuid_regex.MatchString(update.UUID)
	if !valid {
		return
	}

	for _, feed := range update.Feeds {
		valid = url_regex.MatchString(feed.Xml)
		if valid {
			valid = url_regex.MatchString(feed.Web)
		}
		if !valid {
			break
		}
	}
	return
}

func opmlFeed(w http.ResponseWriter, r *http.Request) {
}

func generateShortUrl() (shortid string, err error) {
	r := redis.New("", OPML_REDIS_DB, "")
	var resp []byte
	for {
		shortid = shorten.Shorten()
		if len(shortid) == 0 {
			continue
		} else if resp, err = r.Get("OF_" + shortid); err != nil {
			log.Printf("[!] redis error: %s\n", err.Error())
			shortid = ""
			break
		} else if len(resp) == 0 {
			break
		}
	}
	return
}

func clientUpdate(update *Update) (shortid string, err error) {
	r := redis.New("", OPML_REDIS_DB, "")
	id := "OF_" + update.UUID

	var jsonData []byte
	var bshortid []byte
        var clientData Response
        clientData.Feeds = update.Feeds
	if bshortid, err = r.Get("OF_id_" + update.UUID); err != nil {
		return
	} else if len(bshortid) == 0 {
		shortid, err = generateShortUrl()
	} else {
		shortid = string(bshortid)
	}

	if err != nil {
		return
	}

	jsonData, err = json.Marshal(clientData)
        fmt.Println("[+] json data: ", string(jsonData))
	if err != nil {
		return
	}
	if err = r.Set(id, jsonData); err != nil {
		return
	} else if err = r.Set("OF_"+shortid, update.UUID); err != nil {
		return
	} else if err = r.Set("OF_id_"+update.UUID, shortid); err != nil {
		return
	}
	return
}

func shortIdToUUID(shortid string) (uuid string, err error) {
	key := "OF_" + shortid
	r := redis.New("", OPML_REDIS_DB, "")
	bval, err := r.Get(key)
	if err != nil {
		log.Println("[!] redis error: ", err.Error())
	} else if len(bval) > 0 {
		uuid = string(bval)
	}
	return
}

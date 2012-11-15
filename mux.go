package main

import (
	"bitbucket.org/kisom/opmlfeed/shorten"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// LoadUpdate takes incoming JSON data and unpacks it into an update value
func LoadUpdate(jsonData []byte) (update *Update, err error) {
	update = new(Update)
	err = json.Unmarshal(jsonData, &update)
	return
}

// initialise the code necessary for actually serving client responses
// and interacting with the database.
func InitMux() {
	var err error
	index_text, err = ioutil.ReadFile(index_file)
	if err != nil {
		log.Panic("[!] unable to load index file: ", err.Error())
	}
	opml_mux.HandleFunc("/", OpmlRoot)
	initDatabase()
}

// OpmlRoot is the primary routing construct
func OpmlRoot(w http.ResponseWriter, r *http.Request) {
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
		update, err := LoadUpdate(rawJson)
		if err != nil || !update.Validate() {
			res := http.StatusBadRequest
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
		shortid, err := ClientUpdate(update)
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
			uuid, err := uuidFromShort(shortid)
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
			update, err := FetchFeed(string(uuid))
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

/*
 * FetchFeed looks up the given UUID, and returns any associated client
 * data
 */
func FetchFeed(uuid string) (update *Update, err error) {
	opml, err := opmlFromUUID(uuid)
	if len(opml) == 0 {
		return
	}
	if err == nil {
		err = json.Unmarshal(opml, &update)
	}
	return
}

// ShortIdUnused looks up to see if the short code is presently unused
func ShortIdUnused(shortid string) (valid bool, err error) {
	var resp []byte
	resp, err = uuidFromShort(shortid)
	if err != nil || len(resp) > 0 {
		valid = false
	} else {
		valid = true
	}
	return
}

// GenerateShortUrl generates a new short code for a URL.
func GenerateShortUrl() (shortid string, err error) {
	var resp []byte
	for {
		shortid = shorten.Shorten()
		if len(shortid) == 0 {
			continue
		} else if resp, err = uuidFromShort(shortid); err != nil {
			log.Printf("[!] db error: %s\n", err.Error())
			shortid = ""
			break
		} else if len(resp) == 0 {
			break
		}
	}
	return
}

/*
 * ClientUpdate stores the given update to the database, and returns
 * 
 */
func ClientUpdate(update *Update) (shortid string, err error) {
	var jsonData []byte
	var bshortid []byte
	var clientData Response
	clientData.Feeds = update.Feeds
	if bshortid, err = shortFromUUID(update.UUID); err != nil {
		return
	} else if len(bshortid) == 0 {
		shortid, err = shorten.ShortenUrl(ShortIdUnused)
	} else {
		shortid = string(bshortid)
	}

	if err != nil {
		return
	}

	jsonData, err = json.Marshal(clientData)
	if err != nil {
		return
	}
	if err = storeClientData(update.UUID, jsonData); err != nil {
		return
	}
	if err = associateUUIDandShortid(shortid, update.UUID); err != nil {
		return
	}
	return
}

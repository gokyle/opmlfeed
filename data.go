package main

import (
	"net/http"
	"regexp"
)

const (
	OPMLFEED_VERSION = "0.9.3"
	index_file       = "index.txt"
)

// Server configuration variables.
var (
	OPMLFEED_SSL_KEY  string
	OPMLFEED_SSL_CERT string
	SERVER_ADDR       string
	SERVER_PORT       string
)

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

// type Feed represents a single feed
type Feed struct {
	Title string `json:"title"`
	Xml   string `json:"xml"`
	Web   string `json:"web"`
}

// type Update represents a client's data as uploaded to the server
type Update struct {
	UUID  string `json:"uuid"`
	Feeds []Feed `json:"feeds"`
}

/*
 * Validate is used to perform basic validation on the update objects;
 * this ensures incoming data from the user is safe. This employs basic
 * validation of URLs via regular expressions.
 */
func (update *Update) Validate() (valid bool) {
	valid = uuid_regex.MatchString(update.UUID)
	if !valid {
		return
	}

	for _, feed := range update.Feeds {
		valid = url_regex.MatchString(feed.Xml)
		if valid {
			if len(feed.Web) > 0 {
				valid = url_regex.MatchString(feed.Web)
			}
		}
		if !valid {
			break
		}
	}
	return
}

/*
 * Response represents the data that should be sent to a client when they
 * request a short code.
 */
type Response struct {
	Feeds []Feed `json:"feeds"`
}

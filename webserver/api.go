package webserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jech/galene/stats"
)

func parseContentType(ctype string) string {
	return strings.Trim(strings.Split(ctype, ";")[0], " ")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		failAuthentication(w, "galene-api")
		return
	}

	if ok, err := adminMatch(username, password); !ok {
		if err != nil {
			log.Printf("Administrator password: %v", err)
		}
		failAuthentication(w, "galene-api")
		return
	}

	if !strings.HasPrefix(r.URL.Path, "/galene-api/") {
		http.NotFound(w, r)
		return
	}

	first, kind, rest := splitPath(r.URL.Path[len("/galene/api"):])
	if first != "" {
		http.NotFound(w, r)
		return
	}

	if kind == ".stats" && rest == "" {
		if r.Method != "HEAD" && r.Method != "GET" {
			http.Error(w, "method not allowed",
				http.StatusMethodNotAllowed)
		}
		w.Header().Set("content-type", "application/json")
		w.Header().Set("cache-control", "no-cache")
		if r.Method == "HEAD" {
			return
		}

		ss := stats.GetGroups()
		e := json.NewEncoder(w)
		e.Encode(ss)
		return
	}

	http.NotFound(w, r)
	return
}

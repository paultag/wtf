package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"

	"pault.ag/go/wtf"
)

//
func writeJSON(w http.ResponseWriter, data interface{}, code int) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

//
func writeError(w http.ResponseWriter, message string, code int) error {
	return writeJSON(w, map[string][]map[string]string{
		"errors": []map[string]string{
			map[string]string{"message": message},
		},
	}, code)
}

func main() {
	mux := http.NewServeMux()

	db, err := wtf.NewDatabase(os.Args[1])
	if err != nil {
		panic(err)
	}

	wtfEndpoint := "/wtf/"
	mux.HandleFunc(wtfEndpoint, func(w http.ResponseWriter, req *http.Request) {
		acronym := strings.ToLower(req.URL.Path[len(wtfEndpoint):])

		acronyms := wtf.Acronyms{}
		if err := db.Unpack(acronym, &acronyms); err != nil {
			if err == leveldb.ErrNotFound {
				writeError(w, "No acronym found", 404)
				return
			}
			writeError(w, fmt.Sprintf("Error: %s", err.Error()), 500)
			return
		}

		writeJSON(w, acronyms, 200)
	})

	mux.HandleFunc("/slack/", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			writeError(w, "Error parsing form", 500)
			return
		}

		acronym := strings.ToLower(req.Form.Get("text"))

		acronyms := wtf.Acronyms{}
		if err := db.Unpack(acronym, &acronyms); err != nil {
			if err == leveldb.ErrNotFound {
				writeError(w, "No acronym found", 404)
				return
			}
			writeError(w, fmt.Sprintf("Error: %s", err.Error()), 500)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(200)
		w.Write([]byte(acronyms.String()))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		writeError(w, "No such page", 404)
	})

	panic(http.ListenAndServe(":2838", mux))
}

// vim: foldmethod=marker

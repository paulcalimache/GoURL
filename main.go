package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Shortener struct {
	urls map[string]string
}

func (s *Shortener) shortenURL(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.Write([]byte(fmt.Sprint(http.StatusBadRequest) + " - Bad request"))
		return
	}
	shortURL := req.FormValue("short_url")
	if _, exist := s.urls[shortURL]; exist {
		w.Write([]byte(fmt.Sprint(http.StatusForbidden) + " - Short URL path already exist, please use another one."))
		return
	}
	s.urls[shortURL] = req.FormValue("url")
	log.Info().Msgf("URL %s shortened", req.FormValue("url"))
	http.Redirect(w, req, "/", http.StatusMovedPermanently)
}

func (s *Shortener) handleRequest(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	url, exist := s.urls[vars["key"]]
	if !exist {
		log.Info().Msgf("URL %s doesn't exist", vars["key"])
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
		return
	}
	log.Info().Msgf("Redirect to url: %s", url)
	http.Redirect(w, req, url, http.StatusMovedPermanently)
}

func main() {
	s := Shortener{make(map[string]string)}
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./static/")))
	r.HandleFunc("/shorten", s.shortenURL)
	r.HandleFunc("/{key}", s.handleRequest)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

package shortener

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"slices"

	"github.com/gorilla/mux"
	"github.com/paulcalimache/gourl/internal/db"
	"github.com/paulcalimache/gourl/internal/model"
	"github.com/rs/zerolog/log"
)

type Shortener struct {
	db db.Database
}

func NewShortener(db db.Database) *Shortener {
	return &Shortener{
		db: db,
	}
}

func (s *Shortener) Shorten(w http.ResponseWriter, req *http.Request) {
	body, _ := io.ReadAll(req.Body)
	var urlSchema model.UrlSchema
	err := json.Unmarshal(body, &urlSchema)
	if err != nil {
		log.Info().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlSchema.Url, err = normalizeURL(urlSchema.Url)
	if err != nil {
		log.Info().Err(err).Send()
		http.Error(w, "URL provided is not a valid URL", http.StatusBadRequest)
		return
	}
	// Use custom alias or generate it
	if urlSchema.Alias == "" {
		urlSchema.Alias = generateKey(urlSchema.Url)
	}
	if isAliasForbidden(urlSchema.Alias) {
		log.Info().Str("alias", urlSchema.Alias).Msg("Alias forbidden")
		http.Error(w, "Custom alias "+urlSchema.Alias+" forbidden", http.StatusBadRequest)
		return
	}
	err = s.db.CreateShortURL(urlSchema)
	if err != nil {
		log.Info().Err(err).Send()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info().Str("alias", urlSchema.Alias).Str("url", urlSchema.Url).Msg("New url shortened !")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("URL " + urlSchema.Url + " shortened to /" + urlSchema.Alias))
}

func (s *Shortener) Redirect(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	url, err := s.db.GetURL(vars["key"])
	if err != nil {
		log.Err(err).Str("alias", vars["key"]).Send()
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
		return
	}
	log.Info().Str("alias", vars["key"]).Str("url", url).Msg("Successful redirection")
	http.Redirect(w, req, url, http.StatusMovedPermanently)
}

// normalizeURL ensure the given string is an URL
func normalizeURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	if parsedURL.Host == "" {
		return "", errors.New(url.InvalidHostError.Error(""))
	}
	if parsedURL.Scheme == "" {
		inputURL = "https://" + inputURL
	}
	_, err = url.ParseRequestURI(inputURL)
	if err != nil {
		return "", err
	}
	return inputURL, nil
}

func isAliasForbidden(value string) bool {
	notAllowedValues := []string{"shorten", "healthz", "readyz"}
	return slices.Contains(notAllowedValues, value)
}

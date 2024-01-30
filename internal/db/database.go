package db

import (
	"github.com/paulcalimache/gourl/internal/model"
)

type Database interface {
	// CreateShortURL insert the urlSchema given in the database if the alias is not already present.
	CreateShortURL(urlSchema model.UrlSchema) error
	GetURL(alias string) (string, error)
}

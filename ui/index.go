package ui

import (
	"net/http"
)

func (p *Page) Index(res http.ResponseWriter, req *http.Request) Tags {
	RedirectIfNotAuth(res, req)
	return map[string]interface{}{
		"title": "Cerca le tue ricette preferite",
	}
}

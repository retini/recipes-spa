package ui

import "net/http"

func (p *Page) About(res http.ResponseWriter, req *http.Request) Tags {
	return map[string]interface{}{
		"title": "Questa è la pagina about us",
	}
}

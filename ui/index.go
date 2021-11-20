package ui

import "net/http"

// Index ritorna la pagina "index.html" con la lista dei retailers.
func (p *Page) Index(res http.ResponseWriter, req *http.Request) Tags {
	return map[string]interface{}{
		"variabile": "Questa è la homepage",
	}
}

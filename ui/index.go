package ui

import "net/http"

func (p *Page) Index(res http.ResponseWriter, req *http.Request) Tags {
	return map[string]interface{}{
		"variabile": "Questa è la homepage",
	}
}

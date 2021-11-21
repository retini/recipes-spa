package ui

import "net/http"

func (p *Page) Test(res http.ResponseWriter, req *http.Request) Tags {
	return map[string]interface{}{
		"variabile": "Questa è una variabile test",
	}
}

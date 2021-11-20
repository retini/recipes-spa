package ui

import "net/http"

func (p *Page) Recipes(res http.ResponseWriter, req *http.Request) Tags {
	return map[string]interface{}{"example": 5}
}

func (s *Script) Recipes() string {
	return "dati recuperati dall'api"
}

package ui

import "net/http"

func (p *Page) Logout(res http.ResponseWriter, req *http.Request) {
	RedirectIfNotAuth(res, req)
	ck := &http.Cookie{
		Name:   "auth",
		MaxAge: -1,
	}
	http.SetCookie(res, ck)
	res.Header().Set("Cache-Control", "no-cache")
	http.Redirect(res, req, "http://localhost:8080/login.html", http.StatusMovedPermanently)
}

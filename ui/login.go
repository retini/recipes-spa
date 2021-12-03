package ui

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/open2b/scriggo"
)

var loginUser = "user"
var loginPassword = "user"

func (p *Page) Login(res http.ResponseWriter, req *http.Request) {
	fsys := os.DirFS("ui/assets")
	t, err := scriggo.BuildTemplate(fsys, "login.html", nil)
	if err != nil {
		log.Print("building login template, produced following error ", err)
		return
	}
	err = t.Run(res, nil, nil)
	if err != nil {
		log.Print("running login template, produced following error ", err)
		return
	}
}

// newSecureCookie restituisce sia il cookie http da settare, sia l'istanza di
// securecookie che lo ha prodotto. Per questo è utile sia in fase di creazione
// che in fase di decodifica.
func newSecureCookie() (*http.Cookie, *securecookie.SecureCookie, error) {
	// Hash keys should be at least 32 bytes long
	var hashKey = []byte("very-secret")
	// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
	// Shorter keys may weaken the encryption used.
	var blockKey = []byte("a-lot-of-secrets")
	var sc = securecookie.New(hashKey, blockKey)
	cookieValue := true
	encoded, err := sc.Encode("logged-in", cookieValue)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot encode cookie: %s", err)
	}
	cookie := &http.Cookie{
		Name:     "auth",
		Value:    encoded,
		Path:     "/",
		Secure:   false,
		HttpOnly: false,
	}
	return cookie, sc, nil
}

func cookieIsValid(ck *http.Cookie) (bool, error) {
	_, sc, err := newSecureCookie()
	if err != nil {
		return false, err
	}
	var authValue bool
	err = sc.Decode("logged-in", ck.Value, &authValue)
	if err != nil {
		return false, nil
	}
	return authValue, nil
}

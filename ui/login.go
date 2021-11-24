package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/open2b/scriggo"
)

var loginUser = "user"
var loginPassword = "user"

func (p *Page) Login(res http.ResponseWriter, req *http.Request) {

	// Controlla se esiste un cookie auth. Se non esiste, l'utente è autorizzato
	// a vedere la schermata di login.
	cookie, _ := req.Cookie("auth")

	// Controlla se il valore di auth è valido. In tal caso, non ha senso che
	// l'utente possa nuovamente accedere alla schermata di login, quindi lo
	// reindirizza alla index. Altrimenti continua l'esecuzione.
	if cookie != nil {
		if cookieIsValid(cookie) {
			http.Redirect(res, req, "http://localhost:8080/index.html", http.StatusMovedPermanently)
			return
		}
	}

	// Se il metodo della richiesta è di tipo GET, esegui il template del login
	// screen, contenente il form per l'autenticazione.
	if req.Method == "GET" {
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
		return

		// Se il metodo della richiesta è di tipo POST, controlla i campi inviati
		// dal form di autenticazione e verifica se corrispondono con user e
		// password. In tal caso, esegui il redirect a index. Altrimenti dai errore
		// di autenticazione.
	} else if req.Method == "POST" {
		// login conterrà i campi inviati dal form.
		var login struct {
			User     string
			Password string
		}
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Print("reading request body, produced following error ", err)
			return
		}
		err = json.Unmarshal(data, &login)
		if err != nil {
			log.Print("unmarshaling data from request, produced following error ", err)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		if login.User == loginUser && login.Password == loginPassword {
			ck, _ := newSecureCookie()
			http.SetCookie(res, ck)
			io.WriteString(res, "true")
			return
		}
		io.WriteString(res, "false")
		fmt.Println("credenziali errate")
		return
	}

}

// newSecureCookie restituisce sia il cookie http da settare, sia l'istanza di
// securecookie che lo ha prodotto. Per questo è utile sia in fase di creazione
// che in fase di decodifica.
func newSecureCookie() (*http.Cookie, *securecookie.SecureCookie) {
	// Hash keys should be at least 32 bytes long
	var hashKey = []byte("very-secret")
	// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
	// Shorter keys may weaken the encryption used.
	var blockKey = []byte("a-lot-of-secrets")
	var sc = securecookie.New(hashKey, blockKey)
	cookieValue := true
	encoded, err := sc.Encode("logged-in", cookieValue)
	if err != nil {
		log.Print("encoding the cookie produced following error", err)
	}
	cookie := &http.Cookie{
		Name:     "auth",
		Value:    encoded,
		Path:     "/",
		Secure:   false,
		HttpOnly: false,
	}
	return cookie, sc
}

func cookieIsValid(ck *http.Cookie) bool {
	_, sc := newSecureCookie()
	var authValue bool
	err := sc.Decode("logged-in", ck.Value, &authValue)
	if err != nil {
		fmt.Println("decoding cookie produced following errors: ", err)
	}
	fmt.Println(err)
	return authValue
}

func RedirectIfNotAuth(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("auth")
	if err != nil {
		http.Redirect(res, req, "http://localhost:8080/login.html", http.StatusMovedPermanently)
	}
	if cookie != nil {
		if !cookieIsValid(cookie) {
			http.Redirect(res, req, "http://localhost:8080/login.html", http.StatusMovedPermanently)
		}
	}
}

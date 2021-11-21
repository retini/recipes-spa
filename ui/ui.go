package ui

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/native"
)

// Page contiene i metodi utilizzati per costruire le variabili da inserire nei
// template. Risponde alle richieste di file .HTML.
type Page struct {
}

// Script contiene i metodi utilizzati per chiamare le API ed ottenere dati.
// Risponde alle richieste di file .JSON. Può essere chiamato anche dai metodi
// di Page nel processo di costruzione delle variabili.
type Script struct {
}

// Tags è un map da stringa ad interfaccia che contiene le variabili ottenute da
// un metodo di Page. Viene passato come valore al campo native.declarations di
// Scriggo.BuildOptions{}, il quale viene a sua volta passato a
// Scriggo.BuildTemplate come map di valori accessibile all'interno del
// template.
type Tags map[string]interface{}

func UiHandler(res http.ResponseWriter, req *http.Request) {

	// Prendi la URL della richiesta e puliscila da eventuali imprecisioni.
	var urlPath = req.URL.Path
	urlPath = path.Clean(urlPath)

	// Se l'ultimo carattere della URL è lo slash, allora utilizza lo stesso
	// metodo associato alle richieste per il file index.html.
	if rune(urlPath[len(urlPath)-1]) == '/' {
		urlPath += "index.html"
	}

	// Controlla se la richiesta è arrivata da javascript (il tal caso possiede
	// il cookie from-js).
	isFromJs := false
	for _, cookie := range req.Cookies() {
		if cookie.Name == "is-from-js" {
			isFromJs = true
		}
	}

	// Se la richiesta non è arrivata da javascript, allora l'utente non è
	// dentro all'applicazione, quindi va inviato l'intero layout insieme ad
	// initializer.js e non il singolo componente da innestare.
	if !isFromJs {
		fsys := os.DirFS("ui/assets")
		template, err := scriggo.BuildTemplate(fsys, "layout.html", nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		err = template.Run(res, nil, nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	// Prendi la parte finale della URL, ovvero quella corrispondente al file.
	var _, fileName = path.Split(urlPath)

	// Se il suffisso del file richiesto è HTML, allora chiama il metodo di Page
	// associato al file richiesto dalla URL. Il metodo restituirà i tags
	// contenenti le variabili da usare nel template scriggo.
	if strings.HasSuffix(fileName, ".html") {

		// Instanzia il tipo Page (con tutti i suoi metodi associati).
		page := &Page{}

		// Ottieni il metodo di Page vero e proprio (di fatto, restituisce una
		// funzione). Il metodo è quello che ha lo stesso nome di filename.
		var method = reflect.ValueOf(page).MethodByName(pathToMethod(fileName))

		// Se non esiste alcun metodo con quel nome, allora restituisci una
		// risposta 404. NOTA: ESSENDO GIà STATO CARICATO IL LAYOUT, SE IL
		// METODO NON ESISTE DOBBIAMO COMUNQUE PASSARGLI UN COMPONENTE CHE DICA
		// IL MESSAGGIO 404 ALL'UTENTE.
		if !method.IsValid() {
			http.NotFound(res, req)
			return
		}

		// Se invece il metodo esiste, allora chiama il metodo, il quale
		// restituisce le variabili necessarie per costruire il template,
		// sottoforma di tipo Tags. NOTA: il tipo Tags non è altro che un tag da
		// stringa a interfaccia vuota.
		var values = method.Call([]reflect.Value{reflect.ValueOf(res), reflect.ValueOf(req)})

		// Controlla che i valori siano effettivamente di tipo Tags. Se non lo
		// sono restituisci un Internal Server Error. Se lo sono, assegna i
		// valori alla variabile tags.
		var tags Tags
		if len(values) > 0 {
			switch t := values[0].Interface().(type) {
			case Tags:
				tags = t
			default:
				http.Error(res, http.StatusText(500), 500)
			}
		}

		// Se la variabile tags contiene dei valori, costruisci un template
		// Scriggo con quei valori.
		if tags != nil {

			// Indica la directory in cui sono contenuti i templates.
			fsys := os.DirFS("ui/assets")

			// Costruisci un'istanza di scriggo.BuildOptions con i tags.
			opt := scriggo.BuildOptions{
				Globals: native.Declarations{
					"Tags": &tags,
				},
			}

			// Costruisci un template, cercando il file html con lo stesso nome
			// del file richiesto nella URL e passandogli i tags come variabili
			// globali.
			template, err := scriggo.BuildTemplate(fsys, fileName, &opt)
			if err != nil {
				log.Fatal(err)
				return
			}
			err = template.Run(res, nil, nil)
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		// Se invece il suffisso del file richiesto termina con .json, allora
		// chiama il metodo di Script corrispondente per rispondere con del
		// json.
	} else if strings.HasSuffix(fileName, ".json") {

		// Nel dropshipper, il metodo POST era l'unico metodo consentito per
		// invocare un metodo di script dal client. Perché?
		// if req.Method != "POST" {
		//  http.Error(res, "metodo non consentito", http.StatusMethodNotAllowed)
		//  return
		// }

		// Instanzia il tipo Script (con tutti i metodi associati).
		var script = &Script{}

		// Recupera il solo nome del file richiesto, togliendo l'estensione .json.
		var methodName = fileName[0 : len(fileName)-5]

		// Trova il metodo di Script che ha lo stesso nome del file richiesto.
		var method = reflect.ValueOf(script).MethodByName(methodName)

		// Se non esiste alcun metodo che corrisponde al nome del file
		// richiesto, restituisci uno status 404.
		if !method.IsValid() {
			http.NotFound(res, req)
			return
		}
		var methodType = method.Type()

		// Ottieni il corpo json inviato dalla richiesta e lo immagazzina in args.
		var decoder = json.NewDecoder(req.Body)

		var numIn = methodType.NumIn()
		var args = make([]reflect.Value, numIn)
		for i := 0; i < numIn; i++ {
			var in = methodType.In(i)
			var arg = reflect.New(in)
			err := decoder.Decode(arg.Interface())
			if err != nil {
				log.Fatal(err)
				http.Error(res, http.StatusText(500), 500)
				return
			}
			args[i] = arg.Elem()
		}

		// Ottieni i valori dal metodo, chiamandolo con gli argomenti json
		// (args) presenti nella richiesta.
		var values = method.Call(args)
		var resValues = make([]interface{}, len(values))
		for i, value := range values {
			resValues[i] = value.Interface()
		}

		// Trasforma i valori in json, da inviare come risposta al client.
		js, err := json.Marshal(resValues)
		if err != nil {
			log.Fatal(err)
			http.Error(res, http.StatusText(500), 500)
			return
		}
		res.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		res.Write(js)
	} else {
		// Se il file richiesto non termina con .html o .json allora la
		// richiesta non può essere soddisfatta, quindi restituisci uno status
		// 404.
		http.NotFound(res, req)
	}
}

// Funzione di supporto utilizzata per individuare il metodo di Page da
// richiamare per ottenere i Tags.
func pathToMethod(path string) string {
	var s = make([]byte, len(path))
	if 'a' <= path[0] && path[0] <= 'z' {
		s[0] = byte(path[0] - ('a' - 'A'))
	} else {
		s[0] = path[0]
	}
	var j = 1
	for i := 1; i < len(path) && path[i] != '.'; i++ {
		if path[i] == '-' {
			i++
			s[j] = byte(path[i] - ('a' - 'A'))
		} else {
			s[j] = path[i]
		}
		j++
	}
	return string(s[:j])
}

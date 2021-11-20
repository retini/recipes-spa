package main

import (
	"net/http"
	"recipespa/ui"
)

func main() {
	// Inizializza un nuovo multiplexer diverso da quello standard. Si tratta
	// del map che associa le URL ai rispettivi handlers.
	var mux = http.NewServeMux()

	// ui.UiHandler si occuperà di gestire tutte le richieste, chiamando il metodo
	// di Page corrispondente al file richiesto dalla URL.
	mux.HandleFunc("/", ui.UiHandler)

	// Inizializza un server che ascolta le richieste su localhost:8080, ed ha
	// al suo interno il mux personalizzato.
	httpServer := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	// Crea un handler che si occuperà di gestire le richieste ai file statici,
	// andando a cercare i file partendo dalla cartella "ui/assets".
	fs := http.FileServer(http.Dir("ui/assets"))

	// Aggiungi al multiplexer i casi in cui le richieste avvengano a file css o
	// js, facendo gestire tali richieste al fileserver, il quale servirà i file
	// specificati dalla URL direttamente, senza ulteriori steps. Il risultato è
	// che il path della URL delle richieste viene aggiunto a partire dalla
	// cartella ui/assets, per cercare il file richiesto. Esempio --> GET
	// https://<domain.com>/css/base.css sarà servito con il contenuto di
	// <root>/ui/assets + /css/base.css.
	mux.Handle("/css/", fs)
	mux.Handle("/js/", fs)

	// Avvia il server
	httpServer.ListenAndServe()
}

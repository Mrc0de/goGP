package main

import (
	gmux "github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func initWeb() (error, *http.Server) {
	// only 443 is used in goBP
	// Port 80 Requests will be redirected
	go http.ListenAndServe(":80", http.HandlerFunc(redirect))

	// Start https
	r := gmux.NewRouter()
	r.HandleFunc("/", Home)
	srv := &http.Server{
		Addr:         "0.0.0.0:443",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// return server, no error
	return nil, srv
}

func Home(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("TEST"))
}

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://geekprojex.com/"
	log.Printf("redirect [%s]\r\nTarget: %s", req.URL.String(), target)
	http.Redirect(w, req, target, http.StatusSeeOther)
}

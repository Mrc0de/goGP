package main

import (
	gmux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var store = sessions.NewCookieStore([]byte("aaaaaaaaaaaaaaaa"), []byte("zzzzzzzzzzzzzzzz")) // this will be discarded

// LoadConfig MUST be called before this
func initWeb() (error, *http.Server) {
	// only 443 is used in goBP
	// Port 80 Requests will be redirected
	store = sessions.NewCookieStore([]byte(settings.SessAuthKey), []byte(settings.SessCryptKey))
	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	// Start https
	r := gmux.NewRouter()
	// routes (web)
	r.HandleFunc("/", Home)
	// routes (webSOCKET)
	r.HandleFunc("/wsConnectTo", wsConnectTo)
	// create server
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

func Home(w http.ResponseWriter, r *http.Request) {
	var (
		data struct {
			Ip       string
			Conf     Config
			loggedIn bool
		}
	)
	data.Conf.WsHost = settings.WsHost
	data.Ip = r.RemoteAddr[0:strings.Index(r.RemoteAddr, ":")]
	//If it's IPV6 address, fixup
	if len(data.Ip) < 8 {
		data.Ip = data.Ip + "_IPv6"
	}
	// Session Setup
	session, err := store.Get(r, "goGPSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   43200,
		HttpOnly: false,
		Secure:   true,
	}
	if session.IsNew {
		log.Printf("New Session Created For %s", r.RemoteAddr)
	} else {
		log.Printf("Existing Session Used For %s", r.RemoteAddr)
	}
	loggedIn := session.Values["loggedIn"]
	if loggedIn != nil && loggedIn != false && loggedIn == true {
		log.Printf("User is LoggedIn? -> %s", strconv.FormatBool(loggedIn.(bool)))
		data.loggedIn = loggedIn.(bool)
	} else {
		log.Printf("User is NOT Logged In!")
		session.Values["loggedIn"] = false
		data.loggedIn = false
	}
	sessErr := session.Save(r, w)
	if sessErr != nil {
		log.Printf("Session Save Error: %s", sessErr)
	}
	// Session Ready, send template

	// Templates
	var files []string
	if runtime.GOOS == "windows" {
		files = append(files, "web\\Home.html", "web\\Base.html")
	} else {
		files = append(files, "web/Home.html", "web/Base.html")
	}
	// parse
	tmpl, err := template.ParseFiles(files[0], files[1])
	if err != nil {
		log.Printf("Template Parse Error: %s\n", err)
		writeLog("Template Parse Error: "+err.Error(), false)
	}
	// Exec
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template Exec Error: %s\n", err)
		writeLog("Template Exec Error: "+err.Error(), false)
	}
	writeLog("[ "+data.Ip+" ] - "+r.Method+"  "+r.RequestURI, true)
}

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://geekprojex.com/"
	log.Printf("redirect [%s]\r\nTarget: %s", req.URL.String(), target)
	http.Redirect(w, req, target, http.StatusSeeOther)
}

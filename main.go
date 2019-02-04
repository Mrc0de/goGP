package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

const configFile string = "./settings.json"

var settings Config

func main() {
	err := loadconfig(configFile)
	if err != nil {
		writeLog("Error Loading Config: "+err.Error(), true)
	} else {
		printConfig(false)
	}
	initErr, srv := initWeb()
	if initErr != nil {
		writeLog("Error Starting WebServer: "+initErr.Error(), true)
		os.Exit(1)
	}
	// Listen in BACKGROUND
	go func() {
		if err := srv.ListenAndServeTLS(settings.CertFile, settings.KeyFile); err != nil {
			log.Println(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	writeLog("WebServer Started...", true)
	// Setup Coinbase Websocket Connector

	// Block until  signal.
	<-c
	// Create a deadline
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish")
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	srv.Shutdown(ctx)
	writeLog("Quitting...", true)
	os.Exit(0)

}

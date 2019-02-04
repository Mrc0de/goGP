package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	WsHost   string `json:"wsHost"`
	DbHost   string `json:"dbHost"`
	DbName   string `json:"dbName"`
	DbUser   string `json:"dbUser"`
	DbPass   string `json:"dbPass"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

func loadconfig(filename string) (err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		writeLog("Error Loading settings..."+err.Error(), true)
		return err
	}
	writeLog("Loading settings...", true)
	d := json.NewDecoder(f)
	err = d.Decode(&settings)
	if err != nil {
		writeLog("Error Decoding settings.json..."+err.Error(), true)
		return err
	}
	//writeLog("Websocket Host: "+settings.WsHost, true)
	//writeLog("Database Host: "+settings.DbHost, true)
	return nil
}

func writeLog(msg string, printStdout bool) {
	msg = msg + "\n"
	f, err := os.OpenFile("./goGP.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Error writing to goGP.log: %s\n", err)
	}
	defer f.Close()
	if _, err = f.WriteString("[LOG] - " + msg); err != nil {
		log.Printf("Error writing to goGP.log: %s\n", err)
	}
	if printStdout {
		log.Printf("%s", msg)
	}
}

func printConfig(justLog bool) {
	writeLog("--- Settings ---", !justLog)
	writeLog("wsHost: "+settings.WsHost, !justLog)
	writeLog("dbHost: "+settings.DbHost, !justLog)
	writeLog("dbName: "+settings.DbName, !justLog)
	writeLog("dbUser: "+settings.DbUser, !justLog)
	writeLog("dbPass: ********", !justLog)
	writeLog("certFile: "+settings.CertFile, !justLog)
	writeLog("keyFile: "+settings.KeyFile, !justLog)
	writeLog("---------------", !justLog)

}

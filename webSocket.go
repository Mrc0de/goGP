package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var wsConns []*websocket.Conn

func wsConnectTo(w http.ResponseWriter, r *http.Request) {
	// New Inbound Websocket Upgrade
	var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// Create Connection
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Connect Failed
		writeLog("Websocket Upgrade Error: "+err.Error(), true)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		// Connected!
		wsConns = append(wsConns, wsConn)
		wsConnIp := wsConn.RemoteAddr().String()
		writeLog("Websocket Upgraded: "+wsConnIp, true)
	}
}

package main

import (
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	stationController "github.com/blu-ocean/bo-station-orchestrator/controller"
	"github.com/gorilla/websocket"
)

const (
	pongWait = 60 * time.Second
)

var stationId string
var connectionDetails stationController.ConnectionDetails
var u url.URL

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	stationId = "46d17b96-80a2-11ed-ae44-def9580493cb"
	path := "/station/connection/"
	pathWithParam := path + stationId
	username := "Zion5"
	password := "Babylon@123"

	connectionDetails = stationController.StationMap[stationId]

	u = url.URL{Scheme: "ws", Host: "localhost", Path: pathWithParam}
	log.Printf("connecting to %s", u.String())
	var err error

	h := http.Header{"Authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))}}
	connectionDetails.Ws, _, err = websocket.DefaultDialer.Dial(u.String(), h)
	if err != nil {
		log.Fatal("dial:", err)
	}

	connectionDetails.Ws.SetPongHandler(func(string) error {
		connectionDetails.Ws.SetReadDeadline(
			time.Now().Add(pongWait))
		return nil
	})

	defer connectionDetails.Ws.Close()

	done := make(chan struct{})

	defer close(done)
	go connectionDetails.WritePump()
	time.Sleep(10 * time.Second)
	connectionDetails.Ws.SetPongHandler(func(string) error {
		connectionDetails.Ws.SetReadDeadline(
			time.Now().Add(pongWait))
		return nil
	})
}

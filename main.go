package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "api.p2pquake.net:443", "http service address")

type Message struct {
	Area       int     `json:"area"`
	ID         string  `json:"_id"`
	Code       int     `json:"code"`
	Confidence float64 `json:"confidence"`
	Count      int     `json:"count"`
	CreatedAt  string  `json:"created_at"`
	Expire     string  `json:"expire"`
	Hop        int     `json:"hop"`
	StartedAt  string  `json:"started_at"`
	Time       string  `json:"time"`
	Type       string  `json:"type"`
	UID        string  `json:"uid"`
	UpdatedAt  string  `json:"updated_at"`
	UserAgent  string  `json:"user-agent"`
	Ver        string  `json:"ver"`
	// TODO: "area_confidences" 9611_quake_detect_analysis.json as map
	Areas []struct {
		ID   int `json:"id"`
		Peer int `json:"peer"`
	} `json:"areas"`
	Earthquake []struct {
		DomesticTsunami string `json:"domesticTsunami"`
		ForeignTsunami  string `json:"foreignTsunami"`
		Hypocenter      []struct {
			Name      string  `json:"name"`
			Depth     int     `json:"depth"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Magnitude float64 `json:"magnitude"`
		} `json:"hypocenter"`
		MaxScale int    `json:"maxScale"`
		Time     string `json:"time"`
	} `json:"earthquake"`
	Issue []struct {
		Correct string `json:"correct"`
		Source  string `json:"source"`
		Time    string `json:"time"`
		Type    string `json:"type"`
	} `json:"issue"`
	// TODO: "points" as map (551_quake.json)

}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/v2/ws"}
	log.Printf("connecting to %s", u.String())

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err == websocket.ErrBadHandshake {
		log.Printf("handshake failed with status %d", resp.StatusCode)
	}
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	log.Printf("connected")

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			msg, _ := json.Marshal(message)
			log.Printf("recv: %s", string(msg))
			//log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

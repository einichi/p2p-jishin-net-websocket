package p2pJishinNetWebsocket

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type MessageQuake551 struct {
	Code       int    `json:"code"`
	Id         string `json:"id"`
	Time       string `json:"time"`
	UserAgent  string `json:"user-agent"`
	Version    string `json:"ver"`
	Earthquake struct {
		DomesticTsunami string `json:"domesticTsunami"`
		ForeignTsunami  string `json:"foreignTsunami"`
		Hypocenter      struct {
			Name      string  `json:"name"`
			Depth     int     `json:"depth"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Magnitude float64 `json:"magnitude"`
		} `json:"hypocenter"`
		MaxScale int    `json:"maxScale"`
		Time     string `json:"time"`
	} `json:"earthquake"`
	Issue struct {
		Correct string `json:"correct"`
		Source  string `json:"source"`
		Time    string `json:"time"`
		Type    string `json:"type"`
	} `json:"issue"`
	Points []struct {
		Point struct {
			Address string `json:"addr"`
			IsArea  bool   `json:"isArea"`
			Pref    string `json:"pref"`
			Scale   int    `json:"scale"`
		}
	} `json:"points"`
}

type MessageTsunami552 struct {
	// Not tested
	Code      int    `json:"code"`
	Id        string `json:"id"`
	Time      string `json:"time"`
	Cancelled bool   `json:"cancelled"`
	Issue     struct {
		Source string `json:"source"`
		Time   string `json:"time"`
		Type   string `json:"type"`
	}
	Areas []struct {
		Name      string `json:"name"`
		Grade     string `json:"grade"`
		Immediate bool   `json:"immediate"`
	} `json:"areas"`
}

type MessageEmergencyQuake554 struct {
	Code      int    `json:"code"`
	Id        string `json:"id"`
	Time      string `json:"time"`
	QuakeType string `json:"type"`
}

type MessagePeers555 struct {
	Code      int    `json:"code"`
	CreatedAt string `json:"created_at"`
	Hop       int    `json:"hop"`
	Id        string `json:"id"`
	Time      string `json:"time"`
	Uid       string `json:"uid"`
	Version   string `json:"ver"`
	Areas     []struct {
		ID   int `json:"id"`
		Peer int `json:"peer"`
	} `json:"areas"`
}

type MessageQuakeDetect561 struct {
	Area      int    `json:"area"`
	Code      int    `json:"code"`
	CreatedAt string `json:"created_at"`
	Hop       int    `json:"hop"`
	Id        string `json:"id"`
	Time      string `json:"time"`
	Uid       string `json:"uid"`
	Version   string `json:"ver"`
}

type MessageQuakeDetectAnalysis9611 struct {
	AreaConfidences map[int]AreaConfidencesInner `json:"area_confidences"`
	Code            int                          `json:"code"`
	Confidence      float64                      `json:"confidence"`
	Count           int                          `json:"count"`
	Id              string                       `json:"id"`
	StartedAt       string                       `json:"started_at"`
	Time            string                       `json:"time"`
	UpdatedAt       string                       `json:"updated_at"`
	UserAgent       string                       `json:"user-agent"`
	Version         string                       `json:"ver"`
}

type AreaConfidencesInner struct {
	Confidence float64 `json:"confidence"`
	Count      int     `json:"count"`
	Display    string  `json:"display"`
}

func Subscribe(messageType ...int) {

}

func main() {

	var addr = flag.String("addr", "api.p2pquake.net:443", "http service address")
	logger := log.New(os.Stdout, "[P2P-JISHIN-NET] ", log.Ldate|log.Ltime)
	flag.Parse()

	// TODO: Handle the subscription to certain messages somehow
	subscribe := []int{
		551,
		554,
		555,
		561,
		9611,
	}

	// Handle interrupt signal (ctrl+c)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Specify secure web socket (https)
	u := url.URL{Scheme: "wss", Host: *addr, Path: "/v2/ws"}
	logger.Printf("Connecting to %s", u.String())

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err == websocket.ErrBadHandshake {
		logger.Printf("Error: handshake failed with status %d", resp.StatusCode)
	}
	if err != nil {
		logger.Fatal("Error: ", err)
	}
	defer c.Close()
	logger.Printf("Connected")

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logger.Println("Read:", err)
				logger.Println("Erroneous message: ", string(message))
				return
			}
			var getCode struct {
				Code int `json:"code"`
			}
			if err := json.Unmarshal([]byte(message), &getCode); err != nil {
				logger.Printf("Error: %s", err)
				break
			}
			logger.Printf("Received message with code %d", getCode.Code)
			switch getCode.Code {
			case 551:
				// Earthquake
				if contains(subscribe, 551) {
					var unmarshalMessageQuake551 MessageQuake551
					if err := json.Unmarshal([]byte(message), &unmarshalMessageQuake551); err != nil {
						logger.Printf("Error: %s", err)
						logger.Printf("JSON Message: %s", message)
						break
					}
					// Call processing function
				}
			case 552:
				// Tsunami
				if contains(subscribe, 552) {
					var unmarshalMessageTsunami552 MessageTsunami552
					if err := json.Unmarshal([]byte(message), &unmarshalMessageTsunami552); err != nil {
						logger.Printf("Error: %s", err)
						logger.Printf("JSON Message: %s", message)
						break
					}
					// Debugging, remove this later
					logger.Println("Received Tsunami message, dumping data")
					fmt.Printf("%+v\n", unmarshalMessageTsunami552)
					// End of debug, remove this later
					// Call processing function
				}
			case 554:
				// Emergency Quake
				if contains(subscribe, 554) {
					var unmarshalMessageEmergencyQuake554 MessageEmergencyQuake554
					if err := json.Unmarshal([]byte(message), &unmarshalMessageEmergencyQuake554); err != nil {
						logger.Printf("Error: %s", err)
						logger.Printf("JSON Message: %s", message)
						break
					}
					// Call processing function
				}
			case 555:
				// Peers
				if contains(subscribe, 555) {
					var unmarshalMessagePeers555 MessagePeers555
					if err := json.Unmarshal([]byte(message), &unmarshalMessagePeers555); err != nil {
						logger.Printf("Error: %s", err)
						logger.Printf("JSON Message: %s", message)
						break
					}
					// Call processing function
				}
			case 561:
				// Quake Detect
				if contains(subscribe, 561) {
					var unmarshalMessageQuakeDetect561 MessageQuakeDetect561
					if err := json.Unmarshal([]byte(message), &unmarshalMessageQuakeDetect561); err != nil {
						logger.Printf("Error: %s", err)
						logger.Printf("JSON Message: %s", message)
						break
					}
					// Call processing function
				}
			case 9611:
				// Quake Detect Analysis
				if contains(subscribe, 9611) {
					var unmarshalMessageQuakeDetectAnalysis9611 MessageQuakeDetectAnalysis9611
					if err := json.Unmarshal([]byte(message), &unmarshalMessageQuakeDetectAnalysis9611); err != nil {
						logger.Printf("Error: %s", err)
						logger.Printf("JSON Message: %s", message)
						break
					}
					// Call processing function
				}
			}
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
				logger.Println("Write:", err)
				return
			}
		case <-interrupt:
			logger.Println("Interrupted")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Println("Write close:", err)
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

// For checking if message type is subscribed to
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

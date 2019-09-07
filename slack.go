package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	slackChannel    = "CN22NAL3Y"
	slackToken      = os.Getenv("SLACK_TOKEN")
	slackConnection *websocket.Conn
)

// SlackConn represents a websocket connection to Slack
type SlackConn struct {
	ws               *websocket.Conn
	currentMessageID int
}

// Connect creates a connection to Slack's websocket by authenticating via API and then connecting to the URL we are given
func (s *SlackConn) Connect() (err error) {
	if slackToken == "" {
		panic("No Slack token is set")
	}

	resp, err := http.Get("https://slack.com/api/rtm.connect?token=" + slackToken)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type RTMResponse struct {
		Ok   bool   `json:"ok"`
		URL  string `json:"url"`
		Team struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Domain string `json:"domain"`
		} `json:"team"`
		Self struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"self"`
		Error string `json:"error"`
	}

	var response RTMResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	fmt.Printf("%+v\n", response)
	if response.Ok != true {
		panic("rtm.connect returned an error")
	}

	headers := make(http.Header)
	s.ws, _, err = websocket.DefaultDialer.Dial(response.URL, headers)
	if err != nil {
		return err
	}

	err = nil
	// In the background wait for messages
	go func() {
		for {
			// Block for a message on the socket
			message, err := s.read()
			if err != nil {
				break
			}

			// Decode the json into a generic struct
			var v interface{}
			json.Unmarshal(message, &v)
			data := v.(map[string]interface{})

			switch data["type"] {
			case "hello":
				fmt.Println("Successfully connected to Slack")
				s.sendPing()
			case "message":
				type IncomingMessage struct {
					Type string `json:"type"`
					Ts   string `json:"ts"`
					User string `json:"user"`
					Text string `json:"text"`
				}
				var m IncomingMessage
				json.Unmarshal(message, &m)
				fmt.Println("Got message: " + m.Text)
			default:
				fmt.Printf("%+v\n", data)
			}
		}
	}()

	return err
}

// Read a message off the websocket
func (s *SlackConn) read() (message []byte, err error) {
	_, message, err = s.ws.ReadMessage()
	if err != nil {
		return nil, err
	}

	return message, nil
}

// Write a message to the websocket
func (s *SlackConn) write(data []byte) (err error) {
	err = s.ws.WriteMessage(websocket.TextMessage, data)
	return err
}

// Close tells the remote side we want to disconnect
func (s *SlackConn) Close() (err error) {
	err = s.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return err
}

func (s *SlackConn) sendPing() (err error) {
	type Ping struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
		Time int64  `json:"time"`
	}

	s.currentMessageID++

	var ping Ping
	ping.ID = s.currentMessageID
	ping.Type = "ping"
	ping.Time = time.Now().Unix()

	pingJSON, err := json.Marshal(ping)
	if err != nil {
		return err
	}

	s.write(pingJSON)

	return nil
}

// SendMessage sends the text over the websocket to Slack as a `message` type
func (s *SlackConn) SendMessage(text string) (err error) {
	type OutGoingMessage struct {
		ID      int    `json:"id"`
		Type    string `json:"type"`
		Channel string `json:"channel"`
		Text    string `json:"text"`
	}

	s.currentMessageID++

	var message OutGoingMessage
	message.ID = s.currentMessageID
	message.Type = "message"
	message.Channel = slackChannel
	message.Text = text

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	s.write(messageJSON)

	return nil
}

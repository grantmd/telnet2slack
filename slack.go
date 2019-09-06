package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var (
	slackChannel    = "CN22NAL3Y"
	slackToken      = os.Getenv("SLACK_TOKEN")
	slackConnection *websocket.Conn
)

// SlackConn represents a websocket connection to Slack
type SlackConn struct {
	ws *websocket.Conn
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

	message, err := s.read()
	if err != nil {
		return err
	}

	fmt.Println(string(message))

	return nil
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

// Tell the remote side we want to close
func (s *SlackConn) close() (err error) {
	err = s.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return err
}

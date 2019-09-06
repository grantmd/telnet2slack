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

func slackConnect() {
	if slackToken == "" {
		panic("No Slack token is set")
	}

	resp, err := http.Get("https://slack.com/api/rtm.connect?token=" + slackToken)
	if err != nil {
		fmt.Println("Connect error:")
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
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
	}

	var response RTMResponse
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	fmt.Println(response)
	if response.Ok != true {
		panic("rtm.connect returned an error")
	}

	headers := make(http.Header)
	slackConnection, _, err = websocket.DefaultDialer.Dial(response.URL, headers)
	if err != nil {
		panic(err)
	}

	_, message, err := slackConnection.ReadMessage()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(message))
}

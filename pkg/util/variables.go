// Package constants contains sensitive informations like the serverID and BotToken
package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var ServerID = "9"
var BotToken = "M"

type Config struct {
	BotToken string `json:"BK"`
	ServerID string `json:"SK"`
}

func GetKeys() {
	resp, err := http.Get("https://somedomain.abc/cfg.json")
	if err != nil {
		fmt.Println("Could not establish connection")
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Could not establish connection")
		os.Exit(1)
	}
	var conf Config
	bodyBytes, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &conf)
	if err != nil {
		fmt.Println("Could not parse json correctly")
	}
	ServerID = conf.ServerID
	BotToken = conf.BotToken
}

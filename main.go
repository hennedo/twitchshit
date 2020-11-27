package main

import (
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/sirupsen/logrus"
)

// ENV vars aus secrets.sh: TWITCH_OAUTH_TOKEN_HENNEDO, anwendung: TWITCH_OAUTH_CLIENT_SECRET, TWITCH_OAUTH_CLIENT_ID, erweiterung: TWITCH_API_CLIENT_ID

func main() {


	client := twitch.NewAnonymousClient()

	err := client.Connect()
	if err != nil {
		logrus.Fatal(err)
	}
}

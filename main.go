package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/hennedo/twitchshit/websockets"
	"strings"
)

// ENV vars aus secrets.sh: TWITCH_OAUTH_TOKEN_HENNEDO, anwendung: TWITCH_OAUTH_CLIENT_SECRET, TWITCH_OAUTH_CLIENT_ID, erweiterung: TWITCH_API_CLIENT_ID

var hub *websockets.Hub

func main() {

	flag.Int("port", 8000, "Port where the shop listens on")
	flag.String("host", "", "IP to bind to")
	flag.String("twitch-oauth-token-hennedo", "", "")
	flag.String("twitch-channel", "", "")
	_ = viper.BindPFlags(flag.CommandLine)
	flag.Parse()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	hub = websockets.NewHub()

	hub.On("clap", func(client *websockets.Client, s string) {
		logrus.Info("clapping")
		hub.BroadcastJSON("clap", nil)
	})

	r := mux.NewRouter()
	r.HandleFunc("/ws", hub.ServeWs)

	client := twitch.NewClient(viper.GetString("twitch-channel"), viper.GetString("twitch-oauth-token-hennedo"))
	client.Join(viper.GetString("twitch-channel"))
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if message.Message == "!clap" {
			hub.BroadcastJSON("clap", nil)
		}
	})
	go client.Connect()


	logrus.Error(http.ListenAndServe(fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port")), r))
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gorilla/mux"
	"github.com/hennedo/twitchshit/websockets"
	"github.com/nicklaw5/helix"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

// ENV vars aus secrets.sh: TWITCH_OAUTH_TOKEN_HENNEDO, anwendung: TWITCH_OAUTH_CLIENT_SECRET, TWITCH_OAUTH_CLIENT_ID, erweiterung: TWITCH_API_CLIENT_ID

var hub *websockets.Hub
var helixClient *helix.Client

var twitchToken string

func main() {

	flag.Int("port", 8000, "Port where the shop listens on")
	flag.String("host", "", "IP to bind to")
	flag.String("twitch-oauth-token-hennedo", "", "")
	flag.String("twitch-channel", "", "")
	_ = viper.BindPFlags(flag.CommandLine)
	flag.Parse()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Error(err)
	}


	hub = websockets.NewHub()

	hub.On("clap", func(client *websockets.Client, s string) {
		logrus.Info("clapping")
		hub.BroadcastJSON("clap", nil)
	})

	hub.On("twitchOauth", func(client *websockets.Client, s string) {
		helixClient.SetUserAccessToken(s)
		viper.Set("twitchUserToken", s)
		err := viper.WriteConfigAs("./config.yaml")
		if err != nil {
			logrus.Error(err)
		}
	})

	hub.On("connect", func(client *websockets.Client, s string) {
		users, err := helixClient.GetUsers(&helix.UsersParams{
			Logins: []string{"hennedo92"},
		})
		if err != nil {
			logrus.Error(err)
		}
		follows, err := helixClient.GetUsersFollows(&helix.UsersFollowsParams{
			ToID:   users.Data.Users[0].ID, //todo hennes id rausfinden
		})
		if err != nil {
			logrus.Error(err)
		}
		hub.BroadcastJSON("follow", helix.UserFollow{
			FromName:   follows.Data.Follows[0].FromName,
		})
		broadcastFollowCount()
	})

	r := mux.NewRouter()
	r.HandleFunc("/ws", hub.ServeWs)
	r.HandleFunc("/twitchFollowWebhook", followWebhook)
	r.HandleFunc("/oauthrequest", oauthRequest)
	client := twitch.NewClient(viper.GetString("twitch-channel"), fmt.Sprintf("oauth:%s", viper.GetString("twitchUserToken")))
	client.Join(viper.GetString("twitch-channel"))
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if message.Message == "!clap" {
			hub.BroadcastJSON("clap", nil)
		}
	})
	go func() {
		err := client.Connect()
		if err != nil {
			logrus.Errorf("Chatclient: %s", err)
		}
	}()
	go helixApi()


	logrus.Error(http.ListenAndServe(fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port")), r))
}

func helixApi() {
	var err error
	helixClient, err = helix.NewClient(&helix.Options{
		ClientID:     viper.GetString("twitch-oauth-client-id"),
		ClientSecret: viper.GetString("twitch-ouaht-client-secret"),
		RedirectURI: "http://localhost:8080/oauth/twitch",
	})
	if err != nil {
		logrus.Error(err)
	}
	helixClient.SetUserAccessToken(viper.GetString("twitchUserToken"))
	_, err = helixClient.PostWebhookSubscription(&helix.WebhookSubscriptionPayload{
		Mode:         "subscribe",
		Topic:        "https://api.twitch.tv/helix/users/follows?first=1&to_id=422342243",
		Callback:     "https://twitchshit.dre.li/twitchFollowWebhook",
		LeaseSeconds: 0,
		Secret:       "",
	})
	if err != nil {
		logrus.Error(err)
	}
}

type FollowResponse struct {
	Data []helix.UserFollow
}

func followWebhook(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("hub.challenge") != "" {
		w.Write([]byte(r.URL.Query().Get("hub.challenge")))
	}
	var follow FollowResponse
	err := json.NewDecoder(r.Body).Decode(&follow)
	if err != nil {
		logrus.Error(err)
	}
	if len(follow.Data) > 0 {
		hub.BroadcastJSON("follow", follow.Data[0])
	}
	broadcastFollowCount()
}

func oauthRequest(w http.ResponseWriter, r *http.Request) {
	url := helixClient.GetAuthorizationURL(&helix.AuthorizationURLParams{
		ResponseType: "token",
		Scopes:       []string{"chat:read", "chat:edit", "channel:moderate", "whispers:read", "whispers:edit", "channel_editor"},
		State:        "",
		ForceVerify:  false,
	})
	http.Redirect(w, r, url, 302)
}

func broadcastFollowCount() {
	follows, err := helixClient.GetUsersFollows(&helix.UsersFollowsParams{
		ToID:   "422342243",
	})
	if err != nil {
		logrus.Error(err)
	}
	hub.BroadcastJSON("followers", map[string]int{
		"count": follows.Data.Total,
	})
}
package hackyslack

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Const values for cookie access.
const (
	// Cookie key name.
	Cookie = "c"
	Okay   = "Okay"
	Error  = "Error"
)

var (
	commands = map[string]Command{}
	conf     = oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://slack.com/oauth/authorize",
			TokenURL: "https://slack.com/api/oauth.access",
		},
	}
)

// D is the type for JSON struct.
type D map[string]interface{}

// Args is the Slack command/request struct.
type Args struct {
	TeamID      string
	TeamDomain  string
	ChannelID   string
	ChannelName string
	UserId      string
	UserName    string
	Command     string
	Text        string
	ResponseUrl string
}
type Command func(Args) D

// TeamToken is manually inline oauth2 Token for datastore.
type TeamToken struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
	TeamID       string    `json:"team_id,omitempty"`
	TeamName     string    `json:"team_name,omitempty"`
	Scope        string    `json:"scope,omitempty"`
	Created      time.Time `json:"created,omitempty"`
}

// Configure takes slack ClientID and ClientSecret
func Configure(clientID string, clientSecret string) {
	conf.ClientID = clientID
	conf.ClientSecret = clientSecret
}

// Register commands for slack call.
func Register(name string, cmd Command) {
	commands["/"+name] = cmd
}

// writeJSON to http response.
func writeJSON(w http.ResponseWriter, r *http.Request, data D) {
	bytes, err := json.Marshal(data)
	if err != nil {
		c := appengine.NewContext(r)
		log.Errorf(c, "Failed to mashal %v: %v", data, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// Route for http request process.
func Route(w http.ResponseWriter, r *http.Request) {
	args := Args{
		TeamID:      r.FormValue("team_id"),
		TeamDomain:  r.FormValue("team_domain"),
		ChannelID:   r.FormValue("channel_id"),
		ChannelName: r.FormValue("channel_name"),
		UserId:      r.FormValue("user_id"),
		UserName:    r.FormValue("user_name"),
		Command:     r.FormValue("command"),
		Text:        r.FormValue("text"),
		ResponseUrl: r.FormValue("response_url"),
	}
	// c := appengine.NewContext(r)
	// log.Infof(c, "Got command %v", args)
	cmd, ok := commands[args.Command]
	if ok {
		writeJSON(w, r, cmd(args))
	} else {
		w.Write([]byte("Command not found."))
	}
}

// Oauth for register bot to slack.
func Oauth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"]
	if len(code) == 0 {
		http.Redirect(w, r, "/", 303)
		return
	}
	c := appengine.NewContext(r)
	tok, err := conf.Exchange(c, code[0])
	if err != nil || !tok.Valid() {
		log.Errorf(c, "Failed to exchange token %v: %v", tok, err)
		http.SetCookie(w, &http.Cookie{
			Name:  Cookie,
			Value: Error,
		})
		http.Redirect(w, r, "/", 303)
		return
	}
	team := TeamToken{
		AccessToken:  tok.AccessToken,
		TokenType:    tok.TokenType,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry,
		TeamID:       tok.Extra("team_id").(string),
		TeamName:     tok.Extra("team_name").(string),
		Scope:        tok.Extra("scope").(string),
		Created:      time.Now(),
	}
	key := datastore.NewKey(c, "token", team.TeamID, 0, nil)
	datastore.Put(c, key, &team)
	http.SetCookie(w, &http.Cookie{
		Name:  Cookie,
		Value: Okay,
	})
	http.Redirect(w, r, "/", 303)
}

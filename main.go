package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	host = "irc.twitch.tv"
	port = "6667"
)

func checkError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

// TODO:
//	impl addcmd command (store commands in file)

func main() {
	oauthTokenRaw, err := ioutil.ReadFile("./secret/oauth_token.txt")
	checkError(err)
	oauthToken := strings.Trim(string(oauthTokenRaw), " \n")

	b := newBot("gobo_cseea", oauthToken, "#cseea")
	b.connect()
	defer b.disconnect()

	b.addCmd("ping", cmdPing)
	b.addCmd("twitchbot", cmdTwitchBot)
	b.addCmd("echo", cmdEcho)
	b.addCmd("say", cmdSay)
	b.addCmd("me", cmdMe)
	b.addCmd("addcmd", cmdAddCmd)

	b.addCmdsFromDB()

	b.run()
}

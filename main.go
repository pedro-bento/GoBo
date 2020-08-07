package main

import (
	"io/ioutil"
	"strings"
)

const (
	host = "irc.twitch.tv"
	port = "6667"
)

func checkError(e error) {
	if e != nil {
		// fmt.Println(e)
	}
}

// TODO:
//	clean up code :p

// cmd call premisions

// aritmetics
// !+ , !* , !- , !/

// maybe limit the amount of commands in compositions

// REPL

func main() {
	oauthTokenRaw, err := ioutil.ReadFile("./secret/oauth_token.txt")
	checkError(err)
	oauthToken := strings.Trim(string(oauthTokenRaw), " \n")

	b := newBot("gobo_cseea", oauthToken, "#cseea")
	b.connect()
	defer b.disconnect()

	// add builtins
	b.addCmd("echo", cmdEcho)
	b.addCmd("me", cmdMe)
	b.addCmd("insert", cmdInsert)
	b.addCmd("addcmd", cmdAddCmd)   // streamer / bot
	b.addCmd("addrcmd", cmdAddRCmd) // streamer / bot

	b.addCmdsFromDB()

	// add recurrents
	b.addRCmdsFromDB()

	b.run()
}

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
// cmd call premisions

// maybe limit the amount of commands in compositions

// REPL

// ????????    Multi Level Parsing
// !parse 2 !insert !+ %% $ !square %% <<
// !bin !+ (!square) (!id) x

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
	b.addCmd("repeat", cmdRepeat)

	b.addCmd("+", cmdAdd)
	b.addCmd("-", cmdMinus)
	b.addCmd("*", cmdMult)
	b.addCmd("/", cmdDiv)
	b.addCmd("%", cmdMod)

	b.addCmd("addcmd", cmdAddCmd)   // streamer / bot
	b.addCmd("addrcmd", cmdAddRCmd) // streamer / bot

	b.addCmdsFromDB()

	// add recurrents
	b.addRCmdsFromDB()

	b.run()
}

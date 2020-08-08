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
// .cfg to store nickname, channle, oauth_token
// UI
// REPL
// If a cmd already exist do not add it to db
// maybe limit the amount of commands in compositions

func main() {
	oauthTokenRaw, err := ioutil.ReadFile("./secret/oauth_token.txt")
	checkError(err)
	oauthToken := strings.Trim(string(oauthTokenRaw), " \n")

	b := newBot("gobo_cseea", oauthToken, "#cseea")
	b.connect()
	defer b.disconnect()

	// add builtins
	b.addCmd("echo", command{cmdEcho, badgeNone})
	b.addCmd("me", command{cmdMe, badgeNone})
	b.addCmd("insert", command{cmdInsert, badgeNone})

	// maybe just for mods?
	b.addCmd("repeat", command{cmdRepeat, badgeNone})

	b.addCmd("+", command{cmdAdd, badgeNone})
	b.addCmd("-", command{cmdMinus, badgeNone})
	b.addCmd("*", command{cmdMult, badgeNone})
	b.addCmd("/", command{cmdDiv, badgeNone})
	b.addCmd("%", command{cmdMod, badgeNone})

	b.addCmd("addcmd", command{cmdAddCmd, badgeBroadcaster})   // !addcmd PERM NAME BODY
	b.addCmd("addrcmd", command{cmdAddRCmd, badgeBroadcaster}) // !addRcmd NAME DURATION

	b.addCmdsFromDB()

	// add recurrents
	b.addRCmdsFromDB()

	b.run()
}

// ????????    Multi Level Parsing
// !parse 2 !insert !+ %% $ !square %% <<
// !bin !+ (!square) (!id) x

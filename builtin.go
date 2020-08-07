package main

import (
	"strings"
	"time"
)

func cmdEcho(b *bot, pipe bool, args ...string) string {
	arg := strings.Join(args[1:], " ")

	result := arg

	if pipe {
		return result
	}

	b.sendToChat(result)
	return ""
}

func cmdMe(b *bot, pipe bool, args ...string) string {
	user := args[0]
	arg := strings.Join(args[1:], " ")

	result := user + " " + arg

	if pipe {
		return result
	}

	b.sendToChat(result)
	return ""
}

func cmdAddCmd(b *bot, pipe bool, args ...string) string {
	str := strings.Join(args[1:], " ")

	cmdName, cmdBody, cmdFunc := cmdFromString(str)

	b.addCmd(cmdName, cmdFunc)
	b.addCmdToDB(cmdName + " " + cmdBody)

	return ""
}

func cmdAddRCmd(b *bot, pipe bool, args ...string) string {
	argsTemp := strings.Split(strings.TrimSpace(args[1]), " ")
	cmdName := argsTemp[0]
	cmdDelta := argsTemp[1]

	b.addRCmd(cmdName, cmdDelta)
	b.addRCmdToDB(cmdName + " " + cmdDelta)

	dur, err := time.ParseDuration(cmdDelta)
	checkError(err)
	go b.runRecurrent(rcmd{cmdName, dur})

	return ""
}

func cmdInsert(b *bot, pipe bool, args ...string) string {
	arg := strings.Join(args[1:], " ")

	argTemp := strings.Split(arg, "<<")
	arg1 := argTemp[1]
	arg2 := argTemp[0]

	arg2Temp := strings.Split(arg2, "%%")
	arg2a := arg2Temp[0]
	arg2b := arg2Temp[1]

	result := arg2a + " " + arg1 + " " + arg2b

	if pipe {
		return result
	}

	b.sendToChat(result)
	return ""
}

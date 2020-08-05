package main

import "strings"

type command func(*bot, bool, ...string) string

func cmdPing(b *bot, pipe bool, args ...string) string {
	result := "pong"

	if pipe {
		return result
	}

	b.sendToChat(result)
	return ""
}

func cmdTwitchBot(b *bot, pipe bool, args ...string) string {
	result := "no one: absolutely no one: programming streamers: let's make another twitch bot LUL - arn4vv 2020"

	if pipe {
		return result
	}

	b.sendToChat(result)
	return ""
}

func cmdEcho(b *bot, pipe bool, args ...string) string {
	arg := strings.Join(args[1:], " ")

	result := arg

	if pipe {
		return result
	}

	b.sendToChat(result)
	return ""
}

func cmdSay(b *bot, pipe bool, args ...string) string {
	user := args[0]
	arg := strings.Join(args[1:], " ")

	result := user + " said, " + arg

	if pipe {
		return result
	}

	b.sendToChat(result)
	return ""
}

package main

import "strings"

type command func(*bot, bool, ...string) string

func msgToCmdArg(msg string) (string, string) {
	msg = strings.TrimSpace(msg)
	msg = strings.TrimPrefix(msg, "!")
	cmd := strings.Split(msg, " ")[0]
	arg := strings.TrimPrefix(msg, cmd)
	return cmd, arg
}

func (b *bot) resolveMsg(user, msg, carry string, pipe bool) string {
	cmd, arg := msgToCmdArg(msg)

	if f, ok := b.commads[cmd]; ok {
		return f(b, pipe, user, arg, carry)
	}

	return ""
}

func (b *bot) resolveCompose(user string, composition []string) {
	fst := composition[0]
	lst := composition[len(composition)-1]
	it := composition[1 : len(composition)-1]

	// Get the first result applying the last command of the composition to it's argument.
	carry := b.resolveMsg(user, lst, "", true)

	// Traverse all middle commands carrying the result.
	for i := len(it) - 1; i >= 0; i-- {
		carry = b.resolveMsg(user, it[i], carry, true)
	}

	// Apply the first command with arg + result and send it's result to twitch chat
	b.resolveMsg(user, fst, carry, false)
}

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

package main

import (
	"strings"
)

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

func (b *bot) resolveCompose(user string, pipe bool, composition []string) string {
	carry := ""

	if len(composition) > 1 {
		// Get the first result applying the last command of the composition to it's argument.
		lst := composition[len(composition)-1]
		carry = b.resolveMsg(user, lst, "", true)

		it := composition[1 : len(composition)-1]
		// Traverse all middle commands carrying the result.
		for i := len(it) - 1; i >= 0; i-- {
			carry = b.resolveMsg(user, it[i], carry, true)
		}
	}

	// Apply the first command with arg + result and send it's result to twitch chat
	fst := composition[0]
	return b.resolveMsg(user, fst, carry, pipe)
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

func cmdMe(b *bot, pipe bool, args ...string) string {
	user := args[0]

	result := user

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

func cmdAddCmd(b *bot, pipe bool, args ...string) string {
	user := args[0]
	splited := strings.Split(strings.TrimSpace(args[1]), " ")
	cmdName := splited[0]
	cmdBody := strings.Join(splited[1:], " ")
	cmdBody = cmdBody + strings.Join(args[2:], " ")

	b.addCmdToFile(cmdName + " " + cmdBody)

	b.addCmd(cmdName, func(b1 *bot, pipe1 bool, args1 ...string) string {
		cmdArgs := strings.Join(args1[1:], " ")
		cmd := cmdBody + cmdArgs
		composition := strings.Split(cmd, "$")

		result := b1.resolveCompose(user, true, composition)

		if pipe1 {
			return result
		}

		b.sendToChat(result)
		return ""
	})

	return ""
}

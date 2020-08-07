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

func resolveMsg(b *bot, user, msg, carry string, pipe bool) string {
	cmd, arg := msgToCmdArg(msg)

	if f, ok := b.commands[cmd]; ok {
		return f(b, pipe, user, arg, carry)
	}

	return ""
}

func resolveComposition(b *bot, user string, pipe bool, composition []string) string {
	carry := ""
	if len(composition) > 1 {
		// Get the first result applying the last command of the composition to it's argument.
		lst := composition[len(composition)-1]
		carry = resolveMsg(b, user, lst, "", true)

		it := composition[1 : len(composition)-1]
		// Traverse all middle commands carrying the result.
		for i := len(it) - 1; i >= 0; i-- {
			carry = resolveMsg(b, user, it[i], carry, true)
		}
	}

	// Apply the first command with arg + result and send it's result to twitch chat
	fst := composition[0]
	return resolveMsg(b, user, fst, carry, pipe)
}

func cmdFromString(str string) (string, string, command) {
	splited := strings.Split(strings.TrimSpace(str), " ")
	cmdName := splited[0]
	cmdBody := strings.Join(splited[1:], " ")

	return cmdName, cmdBody, func(b *bot, pipe bool, args ...string) string {
		cmdArgs := strings.Join(args[1:], " ")
		cmd := cmdBody + cmdArgs
		composition := strings.Split(cmd, "$")

		result := resolveComposition(b, args[0], true, composition)

		if pipe {
			return result
		}

		b.sendToChat(result)
		return ""
	}
}

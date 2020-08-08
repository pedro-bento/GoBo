package main

import (
	"strings"
	"time"
)

type command struct {
	action func(*bot, bool, ...string) string
	perm   badge
}

type rcmd struct {
	cmdName string
	delta   time.Duration
}

func msgToCmdArg(msg string) (string, string) {
	msg = strings.TrimSpace(msg)
	msg = strings.TrimPrefix(msg, "!")
	cmd := strings.Split(msg, " ")[0]
	arg := strings.TrimPrefix(msg, cmd)
	return cmd, arg
}

func resolveMsg(b *bot, userBadge badge, user, msg, carry string, pipe bool) string {
	cmd, arg := msgToCmdArg(msg)

	if c, ok := b.commands[cmd]; ok {
		if asPerm(c.perm, userBadge) {
			return c.action(b, pipe, user, arg, carry)
		}
		b.sendToChat(user + " " + "you do not have enough permissions to use " + cmd)
	}

	return ""
}

func resolveComposition(b *bot, userBadge badge, user string, pipe bool, composition []string) string {
	carry := ""
	if len(composition) > 1 {
		// Get the first result applying the last command of the composition to it's argument.
		lst := composition[len(composition)-1]
		carry = resolveMsg(b, userBadge, user, lst, "", true)

		it := composition[1 : len(composition)-1]
		// Traverse all middle commands carrying the result.
		for i := len(it) - 1; i >= 0; i-- {
			carry = resolveMsg(b, userBadge, user, it[i], carry, true)
		}
	}

	// Apply the first command with arg + result and send it's result to twitch chat
	fst := composition[0]
	return resolveMsg(b, userBadge, user, fst, carry, pipe)
}

func cmdFromString(str string) (string, string, string, command) {
	splited := strings.Fields(str)
	cmdPermStr := splited[0]
	cmdPerm := badgeFromString(cmdPermStr)
	cmdName := splited[1]
	cmdBody := strings.Join(splited[2:], " ")

	return cmdName, cmdBody, cmdPermStr, command{
		func(b *bot, pipe bool, args ...string) string {
			cmdArgs := strings.Join(args[1:], " ")
			cmd := cmdBody + cmdArgs
			composition := strings.Split(cmd, "$")

			result := resolveComposition(b, cmdPerm, args[0], true, composition)

			if pipe {
				return result
			}

			b.sendToChat(result)
			return ""
		},
		cmdPerm,
	}
}

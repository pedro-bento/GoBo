package main

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"strings"
)

type bot struct {
	nickname   string
	oauthToken string
	channel    string
	conn       net.Conn
	reader     *textproto.Reader
	writer     *textproto.Writer
	commads    map[string]command
}

func newBot(nickname, oauthToken, channel string) bot {
	return bot{
		nickname:   nickname,
		oauthToken: oauthToken,
		channel:    channel,
		commads:    make(map[string]command),
	}
}

func (b *bot) addCmd(cmd string, f command) {
	b.commads[cmd] = f
}

func (b *bot) send(data string) {
	err := b.writer.PrintfLine(data)
	checkError(err)
}

func (b *bot) receive() string {
	data, err := b.reader.ReadLine()
	checkError(err)
	return data
}

func (b *bot) connect() {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	checkError(err)
	b.conn = conn

	reader := bufio.NewReader(b.conn)
	b.reader = textproto.NewReader(reader)

	writer := bufio.NewWriter(b.conn)
	b.writer = textproto.NewWriter(writer)

	b.send("PASS " + b.oauthToken)
	b.send("NICK " + b.nickname)
	b.send("JOIN " + b.channel)

	fmt.Println("Connected.")
}

func (b *bot) disconnect() {
	b.send("PART " + b.channel)
	b.conn.Close()
	fmt.Println("Disconnected.")
}

func (b *bot) sendToChat(msg string) {
	b.send("PRIVMSG " + b.channel + " :" + msg + "\r\n")
}

func (b *bot) handlePing(data string) {
	toSend := strings.ReplaceAll(data, "PING", "PONG")
	fmt.Println("Responding to server PING -> ", toSend)
	b.send(toSend)
}

func msgToCmdArg(msg string) (string, string) {
	msg = strings.TrimSpace(msg)
	msg = strings.TrimPrefix(msg, "!")
	cmd := strings.Split(msg, " ")[0]
	arg := strings.TrimPrefix(msg, cmd)
	return cmd, arg
}

func (b *bot) handlePrivmsg(data string) {
	fields := strings.Split(data, ".tmi.twitch.tv PRIVMSG "+b.channel+" :")

	splitedUserField := strings.Split(fields[0], "@")
	user := splitedUserField[len(splitedUserField)-1]
	msg := fields[1]

	if !strings.HasPrefix(msg, "!") {
		return
	}

	// TODO:
	//	Clean up.
	if msgs := strings.Split(msg, "$"); len(msgs) > 1 {
		fst := msgs[0]
		lst := msgs[len(msgs)-1]
		it := msgs[1 : len(msgs)-1]

		var result string

		// Get the first result applying the last command of the composition to it's argument.
		cmd, arg := msgToCmdArg(lst)
		if f, ok := b.commads[cmd]; ok {
			result = f(b, true, user, arg)
		} else {
			b.sendToChat("@" + user + " unknown command")
			return
		}

		// Traverse all middle commands carrying the result.
		for i := len(it) - 1; i >= 0; i-- {
			cmd, arg := msgToCmdArg(it[i])
			if f, ok := b.commads[cmd]; ok {
				result = f(b, true, user, arg, result)
			} else {
				b.sendToChat("@" + user + " unknown command")
				return
			}
		}

		// Apply the first command with arg + result and send it's result to twitch chat
		cmd, arg = msgToCmdArg(fst)
		if f, ok := b.commads[cmd]; ok {
			result = f(b, false, user, arg, result)
		} else {
			b.sendToChat("@" + user + " unknown command")
			return
		}

	} else {
		cmd, arg := msgToCmdArg(msg)

		if f, ok := b.commads[cmd]; ok {
			f(b, false, user, arg)
		} else {
			b.sendToChat("@" + user + " unknown command")
			return
		}
	}
}

func (b *bot) isQuit() {
	var buff string
	fmt.Scanln(&buff)
	b.disconnect()
	os.Exit(0)
}

func (b *bot) run() {
	for {
		go b.isQuit()

		data := b.receive()
		switch {
		case strings.Contains(data, "PING"):
			go b.handlePing(data)
		case strings.Contains(data, "PRIVMSG"):
			go b.handlePrivmsg(data)
		default:
			{
			}
		}
	}
}

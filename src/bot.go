package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/textproto"
	"os"
	"strings"
	"time"
)

type bot struct {
	nickname          string
	oauthToken        string
	channel           string
	conn              net.Conn
	reader            *textproto.Reader
	writer            *textproto.Writer
	commands          map[string]command
	recurrentCommands []rcmd
	dbCmdsFilepath    string
	dbRCmdsFilepath   string
}

func newBot(nickname, oauthToken, channel string) bot {
	return bot{
		nickname:        nickname,
		oauthToken:      oauthToken,
		channel:         channel,
		commands:        make(map[string]command),
		dbCmdsFilepath:  "./db/cmds.txt",
		dbRCmdsFilepath: "./db/rcmds.txt",
	}
}

func (b *bot) addCmd(cmd string, f command) {
	b.commands[cmd] = f
}

func (b *bot) addCmdsFromDB() {
	data, err := ioutil.ReadFile(b.dbCmdsFilepath)
	checkError(err)
	if len(string(data)) > 0 {
		cmds := strings.Split(strings.TrimSpace(string(data)), "\n")
		for _, cmd := range cmds {
			cmdName, _, _, cmdStruc := cmdFromString(cmd)
			b.addCmd(cmdName, cmdStruc)
		}
	}
}

func (b *bot) addCmdToDB(data string) {
	file, err := os.OpenFile(b.dbCmdsFilepath, os.O_APPEND|os.O_WRONLY, 0777)
	checkError(err)
	defer file.Close()
	_, err = file.WriteString("\n" + data)
	checkError(err)
}

func (b *bot) addRCmd(cmdName, duration string) {
	dur, err := time.ParseDuration(duration)
	checkError(err)
	b.recurrentCommands = append(b.recurrentCommands, rcmd{cmdName, dur})
}

func (b *bot) addRCmdsFromDB() {
	data, err := ioutil.ReadFile(b.dbRCmdsFilepath)
	checkError(err)
	if len(string(data)) > 0 {
		rcmds := strings.Split(strings.TrimSpace(string(data)), "\n")
		for _, rcmd := range rcmds {
			rcmdTemp := strings.Split(rcmd, " ")
			rcmdName := rcmdTemp[0]
			rcmdDelta := rcmdTemp[1]
			b.addRCmd(rcmdName, rcmdDelta)
		}
	}
}

func (b *bot) addRCmdToDB(data string) {
	file, err := os.OpenFile(b.dbRCmdsFilepath, os.O_APPEND|os.O_WRONLY, 0777)
	checkError(err)
	defer file.Close()
	_, err = file.WriteString("\n" + data)
	checkError(err)
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

	b.requestTags()

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

func (b *bot) requestTags() {
	b.send("CAP REQ :twitch.tv/tags")
}

func (b *bot) handlePing(data string) {
	toSend := strings.ReplaceAll(data, "PING", "PONG")
	fmt.Println("Responding to server PING -> ", toSend)
	b.send(toSend)
}

func (b *bot) handlePrivmsg(data string) {
	fields := strings.Split(data, ".tmi.twitch.tv PRIVMSG "+b.channel+" :")
	tagsUserField := strings.Split(fields[0], " :")
	userField := strings.Split(tagsUserField[1], "@")

	tags := tagsUserField[0]
	user := userField[len(userField)-1]
	msg := fields[1]

	// Make sure it's a command.
	if !strings.HasPrefix(msg, "!") {
		return
	}

	userBadge := badgeFromString(tags)

	if strings.Contains(msg, "addcmd") {
		resolveMsg(b, userBadge, user, msg, "", false)
	} else if composition := strings.Split(msg, "$"); len(composition) > 1 {
		resolveComposition(b, userBadge, user, false, composition)
	} else {
		resolveMsg(b, userBadge, user, msg, "", false)
	}
}

func (b *bot) isQuit() {
	var buff string
	fmt.Scanln(&buff)
	b.disconnect()
	os.Exit(0)
}

func (b *bot) runRecurrent(rc rcmd) {
	for {
		if f, ok := b.commands[rc.cmdName]; ok {
			f.action(b, false, b.nickname, "", "")
		}

		time.Sleep(rc.delta)
	}
}

func (b *bot) run() {
	for _, rc := range b.recurrentCommands {
		go b.runRecurrent(rc)
	}

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

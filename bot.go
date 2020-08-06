package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/textproto"
	"os"
	"strings"
)

type bot struct {
	nickname       string
	oauthToken     string
	channel        string
	conn           net.Conn
	reader         *textproto.Reader
	writer         *textproto.Writer
	commads        map[string]command
	dbCmdsFilepath string
}

func newBot(nickname, oauthToken, channel string) bot {
	return bot{
		nickname:       nickname,
		oauthToken:     oauthToken,
		channel:        channel,
		commads:        make(map[string]command),
		dbCmdsFilepath: "./db/cmds.txt",
	}
}

func (b *bot) addCmd(cmd string, f command) {
	b.commads[cmd] = f
}

func (b *bot) cmdFromString(str string) (string, command) {
	splited := strings.Split(strings.TrimSpace(str), " ")
	cmdName := splited[0]
	cmdBody := strings.Join(splited[1:], " ")

	return cmdName, func(b1 *bot, pipe1 bool, args1 ...string) string {
		cmdArgs := strings.Join(args1[1:], " ")
		cmd := cmdBody + cmdArgs
		composition := strings.Split(cmd, "$")

		result := b1.resolveCompose(args1[0], true, composition)

		if pipe1 {
			return result
		}

		b.sendToChat(result)
		return ""
	}
}

func (b *bot) addCmdsFromDB() {
	data, err := ioutil.ReadFile(b.dbCmdsFilepath)
	checkError(err)
	cmds := strings.Split(string(data), "\n")
	for _, cmd := range cmds {
		cmdName, cmdFunc := b.cmdFromString(cmd)
		b.addCmd(cmdName, cmdFunc) 
	}
}

func (b *bot) addCmdToFile(data string) {
	file, err := os.OpenFile(b.dbCmdsFilepath, os.O_APPEND|os.O_WRONLY, 0777)
	checkError(err)
	defer file.Close()
	_, err = file.WriteString(data + "\n")
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

func (b *bot) handlePrivmsg(data string) {
	fields := strings.Split(data, ".tmi.twitch.tv PRIVMSG "+b.channel+" :")
	splitedUserField := strings.Split(fields[0], "@")

	user := splitedUserField[len(splitedUserField)-1]
	msg := fields[1]

	// Make sure it's a command.
	if !strings.HasPrefix(msg, "!") {
		return
	}

	if strings.Contains(msg, "addcmd") {
		b.resolveMsg(user, msg, "", false)
	} else if composition := strings.Split(msg, "$"); len(composition) > 1 {
		b.resolveCompose(user, false, composition)
	} else {
		b.resolveMsg(user, msg, "", false)
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

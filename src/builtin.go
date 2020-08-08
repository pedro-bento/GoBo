package main

import (
	"math"
	"strconv"
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

	cmdName, cmdBody, cmdPerm, cmd := cmdFromString(str)

	b.addCmd(cmdName, cmd)
	b.addCmdToDB(cmdPerm + " " + cmdName + " " + cmdBody)

	return ""
}

func cmdAddRCmd(b *bot, pipe bool, args ...string) string {
	argsTemp := strings.Fields(args[1])
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

	result := strings.Replace(arg2, "%%", arg1, -1)

	if pipe {
		return result
	}

	b.sendToChat(result)
	return ""
}

func cmdRepeat(b *bot, pipe bool, args ...string) string {
	m, _ := strconv.ParseInt(strings.Fields(args[1])[0], 10, 64)
	n := int(m)
	arghead := strings.Join(strings.Fields(args[1])[1:], " ")
	argtail := strings.Join(args[2:], " ")
	arg := strings.TrimSpace(arghead + " " + argtail)

	result := ""
	for i := 0; i < n; i++ {
		result = result + " " + arg
	}

	if pipe {
		return result
	}

	b.sendToChat(result)
	return ""
}

func cmdAdd(b *bot, pipe bool, args ...string) string {
	arg := strings.Join(args[1:], " ")
	numbersStr := strings.Fields(arg)

	result, _ := strconv.ParseFloat(numbersStr[0], 64)
	numbersStr = numbersStr[1:]

	for _, numStr := range numbersStr {
		num, _ := strconv.ParseFloat(numStr, 64)
		result += num
	}

	if pipe {
		return strconv.FormatFloat(result, 'f', 2, 64)
	}

	b.sendToChat(strconv.FormatFloat(result, 'f', 2, 64))
	return ""
}

func cmdMinus(b *bot, pipe bool, args ...string) string {
	arg := strings.Join(args[1:], " ")
	numbersStr := strings.Fields(arg)

	result, _ := strconv.ParseFloat(numbersStr[0], 64)
	numbersStr = numbersStr[1:]

	for _, numStr := range numbersStr {
		num, _ := strconv.ParseFloat(numStr, 64)
		result -= num
	}

	if pipe {
		return strconv.FormatFloat(result, 'f', 2, 64)
	}

	b.sendToChat(strconv.FormatFloat(result, 'f', 2, 64))
	return ""
}

func cmdMult(b *bot, pipe bool, args ...string) string {
	arg := strings.Join(args[1:], " ")
	numbersStr := strings.Fields(arg)

	result, _ := strconv.ParseFloat(numbersStr[0], 64)
	numbersStr = numbersStr[1:]

	for _, numStr := range numbersStr {
		num, _ := strconv.ParseFloat(numStr, 64)
		result *= num
	}

	if pipe {
		return strconv.FormatFloat(result, 'f', 2, 64)
	}

	b.sendToChat(strconv.FormatFloat(result, 'f', 2, 64))
	return ""
}

func cmdDiv(b *bot, pipe bool, args ...string) string {
	arg := strings.Join(args[1:], " ")
	numbersStr := strings.Fields(arg)

	result, _ := strconv.ParseFloat(numbersStr[0], 64)
	numbersStr = numbersStr[1:]

	for _, numStr := range numbersStr {
		num, _ := strconv.ParseFloat(numStr, 64)
		result /= num
	}

	if pipe {
		return strconv.FormatFloat(result, 'f', 2, 64)
	}

	b.sendToChat(strconv.FormatFloat(result, 'f', 2, 64))
	return ""
}

func cmdMod(b *bot, pipe bool, args ...string) string {
	arg := strings.Join(args[1:], " ")
	numbersStr := strings.Fields(arg)

	result, _ := strconv.ParseFloat(numbersStr[0], 64)
	numbersStr = numbersStr[1:]

	for _, numStr := range numbersStr {
		num, _ := strconv.ParseFloat(numStr, 64)
		result = math.Mod(result, num)
	}

	if pipe {
		return strconv.FormatFloat(result, 'f', 2, 64)
	}

	b.sendToChat(strconv.FormatFloat(result, 'f', 2, 64))
	return ""
}

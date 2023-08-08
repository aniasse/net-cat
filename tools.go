package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func ReadUsername(conn net.Conn) string {

	// conn.Write([]byte("[ENTER YOUR NAME]:"))

	name := make([]byte, 1024)
	n, _ := conn.Read(name)
	username := string(name[:n])
	username = strings.Trim(username, "\n")
	username = strings.TrimSpace(username)
	return username
}

func WriteMsg(username, msg string) string {

	if msg == "" {
		times := time.Now().Format("2006-01-02 15:04:05")
		str := fmt.Sprintf("[%s][%v]:", times, username)

		return str
	} else {
		times := time.Now().Format("2006-01-02 15:04:05")
		str := fmt.Sprintf("[%s][%v]:%s", times, username, msg)

		return str
	}
}

func CheckUsername(arrayname []string, username string) bool {

	for _, name := range arrayname {
		if name == username {
			return false
		}
	}

	return true
}

func ChargeNetcat(file string) string {

	content, err := os.ReadFile(file)
	LogError(err)

	return string(content)
}

func Delete(tab []string, str string) []string {

	array := []string{}
	for _, v := range tab {
		if v != str {
			array = append(array, v)
		}
	}
	return array
}

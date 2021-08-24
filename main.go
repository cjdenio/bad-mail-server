package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "25"
	}
	server, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Server started on port %s\n", port)

	client := slack.New(os.Getenv("SLACK_TOKEN"))

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go func(conn net.Conn) {
			// state
			is_data := false
			receiving_data := false

			data := ""

			conn.Write([]byte("220 mail.calebden.io\r\n"))

			scanner := bufio.NewScanner(conn)

			for scanner.Scan() {
				text := scanner.Text()

				if strings.Contains(strings.ToLower(text), "ehlo") {
					conn.Write([]byte("250 mail.calebden.io says howdy\r\n"))
				} else if strings.EqualFold(text, "data") {
					conn.Write([]byte("354 Start mail input; end with <CRLF>.<CRLF>\r\n"))

					is_data = true
				} else if is_data && text == "." {
					conn.Write([]byte("250 OK\r\n"))

					is_data = false
				} else if is_data && (strings.Contains(text, "Subject: ") || strings.Contains(text, "From: ")) {
					data = data + text + "\n"
				} else if is_data && strings.Contains(text, "text/plain") {
					receiving_data = true
				} else if receiving_data && strings.Contains(strings.ToLower(text), "content-type") {
					receiving_data = false

					client.PostMessage("C017MS0S4E6", slack.MsgOptionText(fmt.Sprintf("```%s```", data), false))
				} else if receiving_data {
					data = data + text + "\n"
				} else if !is_data && strings.EqualFold(text, "quit") {
					conn.Write([]byte("221 OK\r\n"))

					conn.Close()
					break
				} else if !is_data {
					conn.Write([]byte("250 OK\r\n"))
				}
			}
		}(conn)
	}
}

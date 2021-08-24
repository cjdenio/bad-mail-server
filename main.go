package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
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

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go func(conn net.Conn) {
			// state
			is_data := false

			fmt.Println("Connected")

			conn.Write([]byte("220 mail.calebden.io\r\n"))

			scanner := bufio.NewScanner(conn)

			for scanner.Scan() {
				fmt.Println("> " + scanner.Text())

				text := scanner.Text()

				if strings.Contains(strings.ToLower(text), "ehlo") {
					conn.Write([]byte("250 mail.calebden.io says howdy\r\n"))
					fmt.Println("< 250 mail.calebden.io says howdy")
				} else if strings.EqualFold(text, "data") {
					conn.Write([]byte("354 Start mail input; end with <CRLF>.<CRLF>\r\n"))
					fmt.Println("< 354 Start mail input; end with <CRLF>.<CRLF>")

					is_data = true
				} else if is_data && text == "." {
					conn.Write([]byte("250 OK\r\n"))
					fmt.Println("< 250 OK")

					is_data = false
				} else if !is_data && strings.EqualFold(text, "quit") {
					conn.Write([]byte("221 OK\r\n"))
					fmt.Println("< 221 OK")

					conn.Close()
					break
				} else if !is_data {
					conn.Write([]byte("250 OK\r\n"))
					fmt.Println("< 250 OK")
				}
			}

			fmt.Println("Disconnected")
		}(conn)
	}
}

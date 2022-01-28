package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

func handler(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	i := 0
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if i == 0 {
			fmt.Println("Method: ", strings.Split(ln, " ")[0])
			fmt.Println("Path: ", strings.Split(ln, " ")[1])
		}
		if ln == "" {
			break
		}
		i++
	}

	body := "I see you connected\n"
	io.WriteString(conn, "HTTP/1.1 200 OK\r\n")
	io.WriteString(conn, "Content-Type: text/html\r\n")
	io.WriteString(conn, "Content-Length: "+strconv.FormatInt(64, len(body))+"\r\n")
	io.WriteString(conn, "\r\n")
	io.WriteString(conn, body)
}

func main() {
	ls, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer ls.Close()

	for {
		conn, err := ls.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go handler(conn)
	}
}

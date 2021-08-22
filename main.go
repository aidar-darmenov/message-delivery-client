package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/aidar-darmenov/message-delivery-client/config"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {

	cfg := config.NewConfiguration("config/config.json")

	tcpAddr, err := net.ResolveTCPAddr("tcp4", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP(cfg.Type, nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	go handleClientIncomingTraffic(conn)
	//handleClientOutgoingTraffic(conn)

	fmt.Println("")
	fmt.Println("Have read message")
}

func handleClientOutgoingTraffic(conn *net.TCPConn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("err reading from input", err)
			return
		}

		sendMessageToServer(conn, text)
	}

}

func sendMessageToServer(conn *net.TCPConn, message string) error {

	var (
		data         = []byte(message)
		msg_len_data = make([]byte, 2)
		buf          = bytes.Buffer{}
	)

	binary.BigEndian.PutUint16(msg_len_data, uint16(len(data)))

	fmt.Println("content length bytes: ", msg_len_data)
	fmt.Println("content bytes: ", data)

	buf.Write(msg_len_data)
	buf.Write(data)

	_, err := conn.Write(buf.Bytes())
	if err != nil {
		return err
	}

	buf.Reset()

	return nil
}

func handleClientIncomingTraffic(conn *net.TCPConn) {
	for {
		var buf [2048]byte
		var contentLength int

		n, err := conn.Read(buf[:2])
		e, ok := err.(net.Error)

		if err != nil && ok && !e.Timeout() {
			fmt.Println(err)
			break
		}

		if n > 0 {
			contentLength = getContentLength(buf[:n])
		} else {
			conn.Write([]byte("n<0"))
		}

		n, err = conn.Read(buf[:contentLength])
		e, ok = err.(net.Error)

		if err != nil && ok && !e.Timeout() {
			fmt.Println(err)
			break
		}

		if n > 0 {
			processContent(buf[:n])
		} else {
			conn.Write([]byte("n<0"))
		}
	}
}

func getContentLength(bufContentLength []byte) int {
	cl := int(binary.BigEndian.Uint16(bufContentLength))
	fmt.Println("content length: ", cl)
	return cl
}

func processContent(buf []byte) {
	fmt.Println("content: " + string(buf))
}

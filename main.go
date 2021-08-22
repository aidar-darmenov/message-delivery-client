package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/aidar-darmenov/message-delivery-client/config"
	"github.com/aidar-darmenov/message-delivery-client/model"
	"log"
	"net"
	"strconv"
)

func main() {

	cfg := config.NewConfiguration("config/config.json")

	var channelMessages chan model.MessageToClients
	channelMessages = make(chan model.MessageToClients, cfg.ChannelMessagesSize)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", net.JoinHostPort(cfg.ConnectionHost, strconv.Itoa(cfg.ConnectionPort)))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP(cfg.ConnectionType, nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	go handleClientIncomingTraffic(conn)
	handleClientOutgoingTraffic(conn, channelMessages)

	fmt.Println("")
	fmt.Println("Client was shut off")
}

func handleClientOutgoingTraffic(conn *net.TCPConn, channelMessages chan model.MessageToClients) {
	for {
		select {
		case message := <-channelMessages:
			err := sendMessageToServer(conn, message)
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}

	}

}

func sendMessageToServer(conn *net.TCPConn, message model.MessageToClients) error {

	var (
		data         []byte
		msg_len_data = make([]byte, 2)
		buf          = bytes.Buffer{}
	)

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint16(msg_len_data, uint16(len(data)))

	fmt.Println("content length bytes: ", msg_len_data)
	fmt.Println("content bytes: ", data)

	buf.Write(msg_len_data)
	buf.Write(data)

	_, err = conn.Write(buf.Bytes())
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

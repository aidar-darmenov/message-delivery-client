package service

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/aidar-darmenov/message-delivery-client/model"
	"net"
)

func (s *Service) SendMessageToClientsByIds(message model.MessageToClients) {
	s.ChannelMessages <- message
}

func (s *Service) SendMessageToServer(conn *net.TCPConn, message model.MessageToClients) error {

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

func (s *Service) HandleClientOutgoingTraffic() {
	for {
		select {
		case message := <-s.ChannelMessages:
			err := s.SendMessageToServer(s.TcpConnection, message)
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
	}
}

func (s *Service) HandleClientIncomingTraffic() {
	for {
		var buf [2048]byte
		var contentLength int

		n, err := s.TcpConnection.Read(buf[:2])
		e, ok := err.(net.Error)

		if err != nil && ok && !e.Timeout() {
			fmt.Println(err)
			break
		}

		if n > 0 {
			contentLength = getContentLength(buf[:n])
		} else {
			s.TcpConnection.Write([]byte("n<0"))
		}

		n, err = s.TcpConnection.Read(buf[:contentLength])
		e, ok = err.(net.Error)

		if err != nil && ok && !e.Timeout() {
			fmt.Println(err)
			break
		}

		if n > 0 {
			processContent(buf[:n])
		} else {
			s.TcpConnection.Write([]byte("n<0"))
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

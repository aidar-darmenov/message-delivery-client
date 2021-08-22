package service

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/aidar-darmenov/message-delivery-client/model"
	"go.uber.org/zap"
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
				s.Logger.Error("Error sending message to server", zap.Error(err))
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
			s.Logger.Error("Error reading content length from TCP connection", zap.Error(err))
			break
		}

		if n > 0 {
			contentLength = s.GetContentLength(buf[:n])
		} else {
			s.TcpConnection.Write([]byte("n<0"))
		}

		n, err = s.TcpConnection.Read(buf[:contentLength])
		e, ok = err.(net.Error)

		if err != nil && ok && !e.Timeout() {
			s.Logger.Error("Error reading content from TCP connection", zap.Error(err))
			break
		}

		if n > 0 {
			s.ProcessContent(buf[:n])
		} else {
			s.TcpConnection.Write([]byte("n<0"))
		}
	}
}

func (s *Service) GetContentLength(bufContentLength []byte) int {
	cl := int(binary.BigEndian.Uint16(bufContentLength))
	s.Logger.Info(fmt.Sprintf("content length: %d", cl))
	return cl
}

func (s *Service) ProcessContent(buf []byte) {
	s.Logger.Info(fmt.Sprintf("content: %v", string(buf)))
	fmt.Println("message from client: ", string(buf))
}

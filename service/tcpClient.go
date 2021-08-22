package service

import (
	"fmt"
)

func (s *Service) StartTcpClient() {

	go s.HandleClientIncomingTraffic()
	s.HandleClientOutgoingTraffic()

	fmt.Println("")
	fmt.Println("Client was shut off")
}

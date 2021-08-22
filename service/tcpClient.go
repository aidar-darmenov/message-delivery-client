package service

import (
	"log"
)

func (s *Service) StartTcpClient() {

	go s.HandleClientIncomingTraffic()
	s.HandleClientOutgoingTraffic()

	log.Fatal("Client was shut off")
}

package service

import (
	"github.com/aidar-darmenov/message-delivery-client/helpers"
	"github.com/aidar-darmenov/message-delivery-client/model"
	"log"
)

func (s *Service) StartTcpClient() {

	s.SendParamsToServer(model.ClientParams{
		HttpPort: s.GetConfigParams().HttpPort,
		Name:     helpers.GenerateRandomString(10),
	})

	go s.HandleClientIncomingTraffic()
	s.HandleClientOutgoingTraffic()

	log.Fatal("Client was shut off")
}

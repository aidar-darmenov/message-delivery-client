package service

import (
	"github.com/aidar-darmenov/message-delivery-client/config"
	"github.com/aidar-darmenov/message-delivery-client/interfaces"
	"github.com/aidar-darmenov/message-delivery-client/model"
	"go.uber.org/zap"
	"log"
	"net"
	"strconv"
)

type Service struct {
	Configuration   interfaces.Configuration
	Logger          *zap.Logger
	ChannelMessages chan model.MessageToClients
	TcpConnection   *net.TCPConn
}

func NewService(cfg *config.Configuration, logger *zap.Logger) *Service {
	//Here can be any other objects like DB, Cache, any kind of delivery services

	tcpAddr, err := net.ResolveTCPAddr("tcp4", net.JoinHostPort(cfg.ConnectionHost, strconv.Itoa(cfg.ConnectionPort)))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP(cfg.ConnectionType, nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	channelMessages := make(chan model.MessageToClients, cfg.ChannelMessagesSize)

	return &Service{
		Configuration:   cfg,
		Logger:          logger,
		ChannelMessages: channelMessages,
		TcpConnection:   conn,
	}
}

func (s *Service) GetLogger() *zap.Logger {
	return s.Logger
}

func (s *Service) GetConfigParams() *config.Configuration {
	return s.Configuration.Params()
}

package interfaces

import (
	"github.com/aidar-darmenov/message-delivery-client/config"
	"github.com/aidar-darmenov/message-delivery-client/model"
	"go.uber.org/zap"
)

type Service interface {
	GetLogger() *zap.Logger
	GetConfigParams() *config.Configuration
	SendMessageToClientsByIds(message model.MessageToClients)
}

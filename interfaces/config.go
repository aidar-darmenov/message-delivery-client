package interfaces

import (
	"github.com/aidar-darmenov/message-delivery-client/config"
)

type Configuration interface {
	InitConfigParams(string)
	Params() *config.Configuration
}

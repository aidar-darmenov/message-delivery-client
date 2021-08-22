package interfaces

import "github.com/gin-gonic/gin"

type WebService interface {
	Start()
	SendMessageToClientsByIds(*gin.Context)
}

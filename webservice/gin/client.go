package gin

import (
	"github.com/aidar-darmenov/message-delivery-client/model"
	"github.com/gin-gonic/gin"
)

func (ws *GinWebService) SendMessageToClientsByIds(c *gin.Context) {

	var message model.MessageToClients

	err := c.Bind(&message)
	if err != nil {
		c.JSON(400, err)
	}

	ws.Service.SendMessageToClientsByIds(message)

	c.JSON(200, "ok")
}

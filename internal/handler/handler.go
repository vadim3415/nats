package handler

import (
	"fmt"

	"nats/internal/model"
	"nats/internal/nats"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var counter int = 0
var cash map[string]model.ModelNats

func (h *Handler) getId(c *gin.Context) {
	id := c.Params.ByName("id")

	if counter == 0 {
		output, err := h.services.PqGetId(id)
		if err != nil {
			logrus.Println(err)
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		cash = make(map[string]model.ModelNats)
		cash[id] = output

		c.JSON(http.StatusOK, cash[id])
		counter++

		return
	}
	fmt.Println("cash")
	c.JSON(http.StatusOK, cash[id])

	return
}

func (h *Handler) natsSub(c *gin.Context) {
	modelNats, err := nats.NatsSubscriber()
	if err != nil {
		logrus.Println(err)
		c.JSON(http.StatusBadRequest, "bad msg")
		return
	}

	err = h.services.PqNatsMsgCreate(modelNats)
	if err != nil {
		logrus.Println(err)
		c.JSON(http.StatusBadRequest, "bad msg")
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"the message from the channel was successfully added to the database": modelNats,
	})

	return
}

func (h *Handler) natsPub(c *gin.Context) {
	err := nats.NatsPublisher()
	if err != nil {
		logrus.Println(err)
		c.JSON(http.StatusOK, "msg not sent")
		return
	}
	c.JSON(http.StatusOK, "sent")

	return
}

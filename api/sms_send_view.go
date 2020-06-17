package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type sendSMSSerializer struct {
	ExternalID string `json:"external_id"`
	Dst        string `json:"dst"`
	Message    string `json:"message"`
	Src        string `json:"src"`
}

func (s *Server) sendSMSView(c *gin.Context) {
	serializer := sendSMSSerializer{}
	if err := c.BindJSON(&serializer); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	data, err := json.Marshal(serializer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "cannot marshal data to json")
		return
	}

	if err := s.Publish(data); err != nil {
		c.JSON(http.StatusInternalServerError, "cannot publish message")
		return
	}
}

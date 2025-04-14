package sse

import (
	"io"
	"notification/internal/storage"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SSEHandler(c *gin.Context) {
	userIDStr := c.Param("userID")

	userIDUint64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid userID"})
		return
	}
	userID := uint32(userIDUint64)

	ch := make(chan string)

	storage.RegisterClient(userID, ch)
	defer storage.UnregisterClient(userID)

	c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-ch:
			c.SSEvent("message", msg)
			return true
		case <-time.After(30 * time.Second):
			c.SSEvent("ping", "keepalive")
			return true
		}
	})
}

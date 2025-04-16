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

	// Установим заголовки, чтобы точно не было проблем
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-ch:
			c.SSEvent("message", msg)
			return true // продолжаем стрим
		case <-time.After(30 * time.Second):
			c.SSEvent("ping", "keepalive")
			return true // продолжаем стрим
		case <-c.Request.Context().Done():
			return false // клиент отключился — завершаем
		}
	})
}

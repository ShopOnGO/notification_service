package sse

import (
	"io"
	"notification/manager"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SSEHandler(cm *manager.ClientManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("userID")
		userIDUint64, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid userID"})
			return
		}
		userID := uint32(userIDUint64)

		ch := make(chan string, 10)
		cm.AddClient(userID, ch)
		defer cm.RemoveClient(userID)

		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin) // или "*" если публичный
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		c.Stream(func(w io.Writer) bool {
			select {
			case msg := <-ch:
				c.SSEvent("message", msg)
				return true
			case <-time.After(30 * time.Second):
				c.SSEvent("ping", "keepalive")
				return true
			case <-c.Request.Context().Done():
				return false
			}
		})
	}
}

func SSEStatusHandler(cm *manager.ClientManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("userID")
		userIDUint64, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid userID"})
			return
		}
		userID := uint32(userIDUint64)

		status := cm.IsConnected(userID)
		c.JSON(200, gin.H{
			"userID":    userID,
			"connected": status,
		})
	}
}

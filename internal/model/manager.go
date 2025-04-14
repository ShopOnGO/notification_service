package model

import (
	"sync"
)

type ClientManager struct {
	Mu      sync.RWMutex
	Clients map[uint32]chan string
}

// Глобальный manager
var Manager = &ClientManager{
	Clients: make(map[uint32]chan string),
}

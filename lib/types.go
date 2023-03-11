package lib

import (
	"net"

	"github.com/google/uuid"
	"github.com/oneforx/go-ecs"
)

type TargetType string

var (
	CLIENT_TARGET TargetType = "CLIENT_TARGET"
	SERVER_TARGET TargetType = "SERVER_TARGET"
)

type Message struct {
	MessageType string      `json:"message"`
	Data        interface{} `json:"data"`
	Target      TargetType
	NetMode     NET_MODE
}

// Connexion Mode
type NET_MODE string

var (
	NET_TCP NET_MODE = "TCP"
	NET_UDP NET_MODE = "UDP"
	NET_HYB NET_MODE = "HYB" // Hybrid, can be sent over the tcp / udp network
)

type Client struct {
	UDP        *net.UDPAddr
	TCP        *net.TCPConn
	Latence    int
	Token      uuid.UUID
	ID         uuid.UUID
	Email      string
	Password   string
	Characters []ecs.Entity
}

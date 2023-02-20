package ecs

import (
	"sync"

	"github.com/google/uuid"
)

type ISystem interface {
	GetName() string
	GetId() uuid.UUID
	Update(*sync.Mutex)
}

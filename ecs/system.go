package ecs

import (
	"github.com/google/uuid"
)

type ISystem interface {
	GetName() string
	GetId() uuid.UUID
	Update()
}

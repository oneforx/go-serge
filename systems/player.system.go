package systems

import (
	"github.com/google/uuid"
	"github.com/oneforx/go-serge/ecs"
)

type PlayerSystem struct {
	Id    uuid.UUID
	Name  string
	World *ecs.IWorld
}

func (ss *PlayerSystem) GetName() string {
	return ss.Name
}

func (ss *PlayerSystem) GetId() uuid.UUID {
	return ss.Id
}

func (ss *PlayerSystem) Update() {
}

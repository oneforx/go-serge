package systems

import (
	"github.com/bawdeveloppement/squirrelgameserver/ecs"
	"github.com/google/uuid"
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

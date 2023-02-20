package library

import (
	"github.com/bawdeveloppement/squirrelgameserver/ecs"
	"github.com/google/uuid"
)

var (
	ENTITY_PLAYER = func(world *ecs.IWorld, id uuid.UUID, ownerId uuid.UUID) *ecs.IEntity {
		var components []*ecs.IComponent = []*ecs.IComponent{
			COMPONENT_POSITION(0, 0),
			COMPONENT_DIMENSION(32, 32),
			COMPONENT_SOLID(true),
		}
		return ecs.CEntityPossessed(world, id, ownerId, components)
	}
	ENTITY_MONSTER = func(world *ecs.IWorld, id uuid.UUID) *ecs.IEntity {
		var components []*ecs.IComponent = []*ecs.IComponent{
			COMPONENT_POSITION(0, 0),
			COMPONENT_DIMENSION(32, 32),
			COMPONENT_SOLID(true),
		}
		return ecs.CEntity(world, id, components)
	}
	ENTITY_BULLET = func(world *ecs.IWorld, id, ownerId uuid.UUID, force, orientation, speed float64) *ecs.IEntity {
		var components []*ecs.IComponent = []*ecs.IComponent{
			COMPONENT_POSITION(0, 0),
			COMPONENT_DIMENSION(32, 32),
			COMPONENT_SOLID(true),
			COMPONENT_FORCE(force),
			COMPONENT_ORIENTATION(orientation),
			COMPONENT_SPEED(speed),
			COMPONENT_DISTANCE(100),
		}
		return ecs.CEntityWithOwner(world, id, ownerId, components)
	}
)

package library

import "github.com/bawdeveloppement/squirrelgameserver/ecs"

var (
	// EN: Component with { x, y, origin_x, origin_y }. OriginXY is used for recognising where the entity was created
	COMPONENT_POSITION = func(x float64, y float64) *ecs.Component {
		return ecs.CreateComponent("position", map[string]interface{}{"x": x, "y": y, "origin_x": x, "origin_y": y})
	}
	COMPONENT_DIMENSION = func(width float64, height float64) *ecs.Component {
		return ecs.CreateComponent("dimension", map[string]interface{}{"width": width, "height": height})
	}
	COMPONENT_SOLID = func(solid bool) *ecs.Component {
		return ecs.CreateComponent("solid", solid)
	}
	COMPONENT_NAME = func(name string) *ecs.Component {
		return ecs.CreateComponent("name", name)
	}
	COMPONENT_SPEED = func(speed float64) *ecs.Component {
		return ecs.CreateComponent("speed", speed)
	}
	COMPONENT_FORCE = func(force float64) *ecs.Component {
		return ecs.CreateComponent("force", force)
	}
	COMPONENT_ORIENTATION = func(orientation float64) *ecs.Component {
		return ecs.CreateComponent("orientation", orientation)
	}
	COMPONENT_DISTANCE = func(distance float64) *ecs.Component {
		return ecs.CreateComponent("distance", distance)
	}
)

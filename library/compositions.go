package library

import "github.com/bawdeveloppement/squirrelgameserver/ecs"

var (
	MOVABLE_COMPOSITION  ecs.Composition = []string{"position"}
	TURNABLE_COMPOSITION ecs.Composition = []string{"orientation"}
	PLAYER_COMPOSITION   ecs.Composition = []string{"position", "dimension", "solid"}
	BULLET_COMPOSITION   ecs.Composition = []string{"position", "dimension", "solid", "speed", "force", "distance", "orientation"}
)

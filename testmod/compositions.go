package main

import "github.com/oneforx/go-serge/ecs"

var (
	MOVABLE_COMPOSITION  ecs.Composition = []string{"oneforx:position"}
	TURNABLE_COMPOSITION ecs.Composition = []string{"oneforx:orientation"}
	PLAYER_COMPOSITION   ecs.Composition = []string{"oneforx:position", "oneforx:dimension", "oneforx:solid"}
	BULLET_COMPOSITION   ecs.Composition = []string{"oneforx:position", "oneforx:dimension", "oneforx:solid", "oneforx:speed", "oneforx:force", "oneforx:distance", "oneforx:orientation"}
)

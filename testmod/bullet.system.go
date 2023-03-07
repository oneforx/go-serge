package main

import (
	"log"

	"github.com/oneforx/go-serge/ecs"
)

type BulletSystem struct {
	Id    ecs.Identifier
	Name  string
	World *ecs.IWorld
}

func (ss *BulletSystem) GetName() string {
	return ss.Name
}

func (ss *BulletSystem) GetId() ecs.Identifier {
	return ss.Id
}

func (ss *BulletSystem) Init(world *ecs.IWorld) {
	ss.World = world
}

func (ss *BulletSystem) Update() {
	worldLocalised := *ss.World
	bulletEntities := worldLocalised.GetEntitiesWithStrictComposition(BULLET_COMPOSITION)
	for _, bullet := range bulletEntities {
		bulletLocalised := *bullet

		bulletPositionComponent := bulletLocalised.GetComponent(ecs.Identifier{
			Namespace: "oneforx",
			Path:      "position",
		})
		bulletPosition, ok := bulletPositionComponent.Data.(map[string]interface{})
		if !ok {
			continue
		}

		x, ok := bulletPosition["x"].(float64)
		if !ok {
			log.Println("dzdazda")
			continue
		}

		bulletPositionComponent.SetData(map[string]interface{}{
			"x":        x + 1,
			"y":        bulletPosition["y"],
			"origin_x": bulletPosition["origin_x"],
			"origin_y": bulletPosition["origin_y"],
		})

		log.Println(bulletPositionComponent.Data)

		// playerEntities := worldLocalised.GetEntitiesWithStrictComposition(library.PLAYER_COMPOSITION)
		// for _, player := range playerEntities {
		// 	playerLocalised := *player
		// 	log.Println(playerLocalised)
		// }
	}
}

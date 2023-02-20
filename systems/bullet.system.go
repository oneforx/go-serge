package systems

import (
	"log"
	"sync"

	"github.com/bawdeveloppement/squirrelgameserver/ecs"
	"github.com/bawdeveloppement/squirrelgameserver/library"
	"github.com/google/uuid"
)

type BulletSystem struct {
	Id    uuid.UUID
	Name  string
	World *ecs.IWorld
}

func (ss *BulletSystem) GetName() string {
	return ss.Name
}

func (ss *BulletSystem) GetId() uuid.UUID {
	return ss.Id
}

func (ss *BulletSystem) Update(sync_mutex *sync.Mutex) {

	sync_mutex.Lock()
	worldLocalised := *ss.World
	bulletEntities := worldLocalised.GetEntitiesWithStrictComposition(library.BULLET_COMPOSITION)
	for _, bullet := range bulletEntities {
		bulletLocalised := *bullet

		bulletPositionComponent := *bulletLocalised.GetComponent("position")
		bulletPosition, ok := bulletPositionComponent.GetData().(map[string]interface{})
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

		// playerEntities := worldLocalised.GetEntitiesWithStrictComposition(library.PLAYER_COMPOSITION)
		// for _, player := range playerEntities {
		// 	playerLocalised := *player
		// 	log.Println(playerLocalised)
		// }
	}
	sync_mutex.Unlock()
}

package game

import (
	"log"
	"sync"

	"github.com/oneforx/go-ecs"
)

type CollisionSystem struct {
	Id    string
	World *ecs.IWorld
}

func (c *CollisionSystem) GetId() string {
	return c.Id
}

func (c *CollisionSystem) Update() {
	var mutex sync.Mutex

	mutex.Lock()
	var world ecs.IWorld = *c.World
	var entities []*ecs.IEntity = world.GetEntitiesWithComponents(ecs.Identifier{Namespace: "oneforx", Path: "position"}, ecs.Identifier{Namespace: "oneforx", Path: "dimension"}, ecs.Identifier{Namespace: "oneforx", Path: "solid"})

	for _, firstEntity := range entities {
		firstEntityLocalised := *firstEntity
		if isInCollision := c.EntityIsInCollisionWith(firstEntity, entities); isInCollision != nil {
			isInCollisionLocalised := *isInCollision
			log.Println("FirstEntity", firstEntityLocalised.GetId().String(), " with ", isInCollisionLocalised.GetId().String())
		}

	}
	mutex.Unlock()
}

func (c *CollisionSystem) EntityIsInCollisionWith(entity *ecs.IEntity, collidablesEntities []*ecs.IEntity) (entityCollided *ecs.IEntity) {
	firstEntityLocalised := *entity
	for _, secondEntity := range collidablesEntities {

		secondEntityLocalised := *secondEntity
		if firstEntityLocalised.GetId() != secondEntityLocalised.GetId() {
			firstPositionComponentPointer := firstEntityLocalised.GetComponent(ecs.Identifier{Namespace: "oneforx", Path: "position"})
			if firstPositionComponentPointer == nil {
				log.Println("skip drawing of : " + firstEntityLocalised.GetId().String())
				continue
			}
			firstPositionComponent := *firstPositionComponentPointer
			firstPositionComponentData, ok := firstPositionComponent.GetData().(map[string]interface{})
			if !ok {
				log.Println("Could not convert firstPositionComponent to map[string]interface{}")
				continue
			}
			firstX := firstPositionComponentData["x"].(int)
			firstY := firstPositionComponentData["y"].(int)

			firstDimensionComponentPointer := firstEntityLocalised.GetComponent(ecs.Identifier{Namespace: "oneforx", Path: "dimension"})
			if firstDimensionComponentPointer == nil {
				log.Println("skip drawing of : " + firstEntityLocalised.GetId().String())
				continue
			}

			firstDimensionComponent := *firstDimensionComponentPointer
			firstDimensionComponentData, ok := firstDimensionComponent.GetData().(map[string]interface{})
			if !ok {
				log.Println("Could not convert firstPositionComponent to map[string]interface{}")
				continue
			}
			firstWidth := firstDimensionComponentData["width"].(int)
			firstHeight := firstDimensionComponentData["height"].(int)

			secondPositionComponentPointer := secondEntityLocalised.GetComponent(ecs.Identifier{Namespace: "oneforx", Path: "position"})
			if secondPositionComponentPointer == nil {
				log.Println("skip drawing of : " + secondEntityLocalised.GetId().String())
				continue
			}
			secondPositionComponent := *secondPositionComponentPointer
			secondPositionComponentData, ok := secondPositionComponent.GetData().(map[string]interface{})
			if !ok {
				log.Println("Could not convert secondPositionComponent to map[string]interface{}")
				continue
			}
			var secondX int = secondPositionComponentData["x"].(int)
			var secondY int = secondPositionComponentData["y"].(int)

			secondDimensionComponentPointer := secondEntityLocalised.GetComponent(ecs.Identifier{Namespace: "oneforx", Path: "dimension"})
			if secondDimensionComponentPointer != nil {
				log.Println("skip drawing of : " + secondEntityLocalised.GetId().String())
				continue
			}
			secondDimensionComponent := *secondDimensionComponentPointer
			secondDimensionComponentData, ok := secondDimensionComponent.GetData().(map[string]interface{})
			if !ok {
				log.Println("Could not convert secondDimensionComponent to map[string]interface{}")
				continue
			}

			var secondWidth int = secondDimensionComponentData["width"].(int)
			var secondHeight int = secondDimensionComponentData["height"].(int)

			if firstX+firstWidth > secondX && firstX < secondX+secondWidth && firstY+firstHeight > secondY && firstY < secondY+secondHeight {
				log.Println(firstEntityLocalised.GetId().String() + " in collision with " + secondEntityLocalised.GetId().String())
			}
		}
	}

	return entityCollided
}

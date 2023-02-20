package ecs

import (
	"errors"

	"github.com/google/uuid"
)

type IEntity interface {
	GetWorld() *IWorld
	GetId() uuid.UUID
	GetOwnerID() uuid.UUID
	GetPossessedID() uuid.UUID
	AddComponent(cmp IComponent) error
	HaveComponent(cn string) bool
	GetComponent(id string) *IComponent
	GetComponents() []*IComponent
	GetComposition() []string
	UpdateComponents([]*IComponent)
	HaveComposition([]string) bool
}

type Entity struct {
	Id          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"ownerId"` // ClientID
	PossessedID uuid.UUID `json:"possessedId"`
	World       *IWorld
	Components  []*IComponent `json:"components"`
}

// Remove Cyclique Structure from type Entity caused by *World qui contient lui même l'entité
type EntityNoCycle struct {
	Id          uuid.UUID     `json:"id"`
	OwnerID     uuid.UUID     `json:"ownerId"`
	PossessedID uuid.UUID     `json:"possessedId"`
	Components  []*IComponent `json:"components"`
}

func (entity *Entity) GetId() uuid.UUID {
	return entity.Id
}

func (entity *Entity) GetOwnerID() uuid.UUID {
	return entity.OwnerID
}

func (entity *Entity) GetPossessedID() uuid.UUID {
	return entity.PossessedID
}

func (entity *Entity) AddComponent(cmp IComponent) error {
	// Check if we already have a component with same id
	var foundId int = -1

	for idx, component := range entity.GetComponents() {
		componentLocalised := *component
		if componentLocalised.GetId() == cmp.GetId() {
			foundId = idx
		}
	}
	if foundId != -1 {
		return errors.New("Component with same id already exist")
	} else {
		entity.Components = append(entity.Components, &cmp)
		return nil
	}
}

// Si un composant spécifié dans l'argument n'existe pas alors on en crée un sur l'entité cible
func (entity *Entity) UpdateComponents(components []*IComponent) {
	for _, c := range components {
		componentLocalised := *c
		targetComponent := entity.GetComponent(componentLocalised.GetId())
		targetComponentLocalised := *targetComponent

		if targetComponentLocalised != nil {
		} else {
			entity.AddComponent(targetComponentLocalised)
		}
	}
}

func (entity *Entity) HaveComponent(cn string) bool {
	for _, component := range entity.Components {
		componentLocalised := *component
		if componentLocalised.GetId() == cn {
			return true
		}
	}
	return false
}

func (entity *Entity) GetComponent(id string) (cmp *IComponent) {
	for _, component := range entity.GetComponents() {
		componentLocalised := *component
		if componentLocalised.GetId() == id {
			cmp = component
		}
	}
	return cmp
}

func (entity *Entity) GetComponents() (components []*IComponent) {
	return entity.Components
}

func (entity *Entity) GetComposition() (composition []string) {
	for _, component := range entity.Components {
		cmp := *component
		composition = append(composition, cmp.GetId())
	}
	return composition
}

func (entity *Entity) HaveComposition(composition []string) bool {
	entityComposition := entity.GetComposition()
	haveComponent := 0
	for _, componentName := range entityComposition {
		for _, targetComponentName := range composition {
			if componentName == targetComponentName {
				haveComponent++
			}
		}
	}
	return len(composition) == haveComponent
}

func (entity *Entity) GetWorld() *IWorld {
	return entity.World
}

func CEntity(world *IWorld, id uuid.UUID, components []*IComponent) *IEntity {
	var newEntity IEntity = &Entity{
		Id:         id,
		World:      world,
		Components: components,
	}
	return &newEntity
}

func CEntityWithOwner(world *IWorld, id uuid.UUID, ownerId uuid.UUID, components []*IComponent) *IEntity {
	var newEntity IEntity = &Entity{
		Id:         id,
		World:      world,
		OwnerID:    ownerId,
		Components: components,
	}
	return &newEntity
}

func CEntityPossessed(world *IWorld, id uuid.UUID, possessedByID uuid.UUID, components []*IComponent) *IEntity {
	var newEntity IEntity = &Entity{
		Id:          id,
		World:       world,
		OwnerID:     possessedByID,
		PossessedID: possessedByID,
		Components:  components,
	}
	return &newEntity
}

package lib

import (
	"errors"

	"github.com/oneforx/go-ecs"
)

// A library system to store content
type ILibrary interface {
	GetId() ecs.Identifier
	GetStruct() Library

	GetCompositions() []ecs.Composition
	GetComponents() []ecs.Component
	GetSystems() []ecs.ISystem

	RegisterSystem(ecs.ISystem) error
	RegisterSystems([]ecs.ISystem) error

	RegisterComponent(ecs.Component) error
	RegisterComponents([]ecs.Component) error
	RegisterComposition(ecs.Composition) error
	RegisterCompositions([]ecs.Composition) error
}
type Library struct {
	Id           ecs.Identifier
	components   []ecs.Component
	systems      []ecs.ISystem
	compositions []ecs.Composition

	serverMessageHandlers map[ecs.Identifier]func(message Message, world *WorldServer)
	clientMessageHandlers map[ecs.Identifier]func(message Message, world *WorldClient)
}

func (lib Library) GetId() ecs.Identifier {
	return lib.Id
}

func (lib Library) GetStruct() Library {
	return lib
}

func (lib Library) GetCompositions() []ecs.Composition {
	return lib.compositions
}

func (lib Library) GetComponents() []ecs.Component {
	return lib.components
}

func (lib Library) GetSystems() []ecs.ISystem {
	return lib.systems
}

func (lib *Library) RegisterSystem(system ecs.ISystem) error {
	if lib.systemExists(system) {
		return errors.New("system already exists")
	}
	lib.systems = append(lib.systems, system)
	return nil
}

func (lib *Library) RegisterSystems(systems []ecs.ISystem) error {
	for _, system := range systems {
		if lib.systemExists(system) {
			return errors.New("system already exists")
		}
	}
	lib.systems = append(lib.systems, systems...)
	return nil
}

func (lib *Library) RegisterComponent(component ecs.Component) error {
	if lib.componentExists(component.Id) {
		return errors.New("component already exists")
	}
	lib.components = append(lib.components, component)
	return nil
}

func (lib *Library) RegisterComponents(components []ecs.Component) error {
	for _, component := range components {
		if lib.componentExists(component.Id) {
			return errors.New("component already exists")
		}
	}
	lib.components = append(lib.components, components...)
	return nil
}

func (lib *Library) RegisterComposition(composition ecs.Composition) error {
	if lib.compositionExists(composition.Id) {
		return errors.New("composition already exists")
	}
	lib.compositions = append(lib.compositions, composition)
	return nil
}

func (lib *Library) RegisterCompositions(compositions []ecs.Composition) error {
	for _, composition := range compositions {
		if lib.compositionExists(composition.Id) {
			return errors.New("composition already exists")
		}
		lib.compositions = append(lib.compositions, composition)
	}
	return nil
}

func (lib *Library) componentExists(id ecs.Identifier) bool {
	for _, component := range lib.components {
		if component.Id == id {
			return true
		}
	}
	return false
}

func (lib *Library) systemExists(system ecs.ISystem) bool {
	for _, s := range lib.systems {
		if s.GetId() == system.GetId() {
			return true
		}
	}
	return false
}

func (lib *Library) compositionExists(id ecs.Identifier) bool {
	for _, composition := range lib.compositions {
		if composition.Id == id {
			return true
		}
	}

	return false
}

func (lib *Library) RegisterHandlerServerMessage(id ecs.Identifier, handler func(Message, *WorldServer)) {
	if lib.Id.Namespace != id.Namespace {
		panic("your message id namespace should be the same as the namespace of your library")
	}

	_, ok := lib.serverMessageHandlers[id]
	if ok {
		panic("un handler avec le même identifiant est déjà enregistrer")
	}

	lib.serverMessageHandlers[id] = handler
}

func (lib *Library) RegisterHandlerClientMessage(id ecs.Identifier, handler func(Message, *WorldClient)) {
	if lib.Id.Namespace != id.Namespace {
		panic("your message id namespace should be the same as the namespace of your library")
	}

	_, ok := lib.clientMessageHandlers[id]
	if ok {
		panic("un handler avec le même identifiant est déjà enregistrer")
	}

	lib.clientMessageHandlers[id] = handler
}

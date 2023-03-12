package goserge

import (
	"errors"

	goecs "github.com/oneforx/go-ecs"
)

// A library system to store content
type ILibrary interface {
	GetId() goecs.Identifier
	GetStruct() Library

	GetCompositions() []goecs.Composition
	GetComponents() []goecs.Component
	GetSystems() []goecs.ISystem

	RegisterSystem(goecs.ISystem) error
	RegisterSystems([]goecs.ISystem) error

	RegisterComponent(goecs.Component) error
	RegisterComponents([]goecs.Component) error
	RegisterComposition(goecs.Composition) error
	RegisterCompositions([]goecs.Composition) error
}
type Library struct {
	Id           goecs.Identifier
	components   []goecs.Component
	systems      []goecs.ISystem
	compositions []goecs.Composition

	serverMessageHandlers map[goecs.Identifier]func(message Message, world *WorldServer)
	clientMessageHandlers map[goecs.Identifier]func(message Message, world *WorldClient)
}

func (lib Library) GetId() goecs.Identifier {
	return lib.Id
}

func (lib Library) GetStruct() Library {
	return lib
}

func (lib Library) GetCompositions() []goecs.Composition {
	return lib.compositions
}

func (lib Library) GetComponents() []goecs.Component {
	return lib.components
}

func (lib Library) GetSystems() []goecs.ISystem {
	return lib.systems
}

func (lib *Library) RegisterSystem(system goecs.ISystem) error {
	if lib.systemExists(system) {
		return errors.New("system already exists")
	}
	lib.systems = append(lib.systems, system)
	return nil
}

func (lib *Library) RegisterSystems(systems []goecs.ISystem) error {
	for _, system := range systems {
		if lib.systemExists(system) {
			return errors.New("system already exists")
		}
	}
	lib.systems = append(lib.systems, systems...)
	return nil
}

func (lib *Library) RegisterComponent(component goecs.Component) error {
	if lib.componentExists(component.Id) {
		return errors.New("component already exists")
	}
	lib.components = append(lib.components, component)
	return nil
}

func (lib *Library) RegisterComponents(components []goecs.Component) error {
	for _, component := range components {
		if lib.componentExists(component.Id) {
			return errors.New("component already exists")
		}
	}
	lib.components = append(lib.components, components...)
	return nil
}

func (lib *Library) RegisterComposition(composition goecs.Composition) error {
	if lib.compositionExists(composition.Id) {
		return errors.New("composition already exists")
	}
	lib.compositions = append(lib.compositions, composition)
	return nil
}

func (lib *Library) RegisterCompositions(compositions []goecs.Composition) error {
	for _, composition := range compositions {
		if lib.compositionExists(composition.Id) {
			return errors.New("composition already exists")
		}
		lib.compositions = append(lib.compositions, composition)
	}
	return nil
}

func (lib *Library) componentExists(id goecs.Identifier) bool {
	for _, component := range lib.components {
		if component.Id == id {
			return true
		}
	}
	return false
}

func (lib *Library) systemExists(system goecs.ISystem) bool {
	for _, s := range lib.systems {
		if s.GetId() == system.GetId() {
			return true
		}
	}
	return false
}

func (lib *Library) compositionExists(id goecs.Identifier) bool {
	for _, composition := range lib.compositions {
		if composition.Id == id {
			return true
		}
	}

	return false
}

func (lib *Library) RegisterHandlerServerMessage(id goecs.Identifier, handler func(Message, *WorldServer)) {
	if lib.Id.Namespace != id.Namespace {
		panic("your message id namespace should be the same as the namespace of your library")
	}

	_, ok := lib.serverMessageHandlers[id]
	if ok {
		panic("un handler avec le même identifiant est déjà enregistrer")
	}

	lib.serverMessageHandlers[id] = handler
}

func (lib *Library) RegisterHandlerClientMessage(id goecs.Identifier, handler func(Message, *WorldClient)) {
	if lib.Id.Namespace != id.Namespace {
		panic("your message id namespace should be the same as the namespace of your library")
	}

	_, ok := lib.clientMessageHandlers[id]
	if ok {
		panic("un handler avec le même identifiant est déjà enregistrer")
	}

	lib.clientMessageHandlers[id] = handler
}

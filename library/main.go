package library

import (
	"fmt"

	"github.com/oneforx/go-serge/ecs"
)

type LibraryManager struct {
	Libraries map[ecs.Identifier]Library
}

func (lm *LibraryManager) GetLibrary(id ecs.Identifier) *Library {
	library, ok := lm.Libraries[id]
	if !ok {
		return nil
	}
	return &library
}

// Return false if we couldn't add the library
func (lm *LibraryManager) AddLibrary(id ecs.Identifier, newLibrary Library) bool {
	if lm.GetLibrary(id) != nil {
		return false
	}

	lm.Libraries[id] = newLibrary
	return true
}

func (lm *LibraryManager) GetComponents() []ecs.Component {
	var components []ecs.Component = []ecs.Component{}

	for _, lib := range lm.Libraries {
		components = append(components, lib.GetComponents()...)
	}

	return components
}

func (lm *LibraryManager) GetComponent(id ecs.Identifier) (cmp ecs.Component) {
	for _, library := range lm.Libraries {
		if library.Id.Namespace == id.Namespace {
			for _, component := range library.GetComponents() {
				if component.Id.Path == id.Path {
					cmp = component
					break
				}
			}
			break
		}
	}
	return cmp
}

func (lm *LibraryManager) GetSystem(id ecs.Identifier) (sys ecs.ISystem, err error) {
	for _, library := range lm.Libraries {
		if library.Id.Namespace == id.Namespace {
			for _, system := range library.GetSystems() {
				if system.GetId().String() == id.String() {
					sys = system
					break
				}
			}
			break
		}
	}

	if sys == nil {
		return nil, fmt.Errorf("System not found")
	}

	return sys, nil
}

func (lm *LibraryManager) InstantiateSystem(id ecs.Identifier, world *ecs.IWorld) (*ecs.ISystem, error) {
	var systemLocation *ecs.ISystem

	system, err := lm.GetSystem(id)
	if err != nil {
		return nil, err
	}

	system.Init(world)

	systemLocation = &system

	return systemLocation, nil
}

func (lm *LibraryManager) InstantiateComponent(id ecs.Identifier, data interface{}) *ecs.Component {
	var componentLocation *ecs.Component

	component := lm.GetComponent(id)

	component.SetData(data)

	componentLocation = &component

	return componentLocation
}

type ILibrary interface {
	GetId() ecs.Identifier
	GetStruct() Library
	GetComponents() []ecs.Component
	GetSystems() []ecs.ISystem
	AddComponent(ecs.Component)
	AddSystem(ecs.ISystem)
}

type Library struct {
	Id         ecs.Identifier
	components []ecs.Component
	systems    []ecs.ISystem
}

func (library *Library) GetId() ecs.Identifier {
	return library.Id
}

func (library *Library) AddSystem(system ecs.ISystem) {
	library.systems = append(library.systems, system)
}

func (library *Library) AddComponent(component ecs.Component) {
	library.components = append(library.components, component)
}

func (library *Library) GetComponents() []ecs.Component {
	return library.components
}

func (library *Library) GetSystems() []ecs.ISystem {
	return library.systems
}

func (library *Library) GetStruct() Library {
	return *library
}

package library

import (
	"fmt"

	"github.com/oneforx/go-ecs"
)

type LibraryManager struct {
	Libraries map[ecs.Identifier]ecs.Library
}

func (lm *LibraryManager) GetLibrary(id ecs.Identifier) *ecs.Library {
	library, ok := lm.Libraries[id]
	if !ok {
		return nil
	}
	return &library
}

// Return false if we couldn't add the library
func (lm *LibraryManager) AddLibrary(id ecs.Identifier, newLibrary ecs.Library) bool {
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

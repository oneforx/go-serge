package lib

import (
	"fmt"

	goecs "github.com/oneforx/go-ecs"
)

type LibraryManager struct {
	Libraries map[goecs.Identifier]ILibrary
}

func (lm *LibraryManager) GetLibrary(id goecs.Identifier) *ILibrary {
	library, ok := lm.Libraries[id]
	if !ok {
		return nil
	}
	return &library
}

// Return false if we couldn't add the library
func (lm *LibraryManager) AddLibrary(id goecs.Identifier, newLibrary ILibrary) bool {
	if lm.GetLibrary(id) != nil {
		return false
	}

	lm.Libraries[id] = newLibrary
	return true
}

func (lm *LibraryManager) GetComponents() []goecs.Component {
	var components []goecs.Component = []goecs.Component{}

	for _, lib := range lm.Libraries {
		components = append(components, lib.GetComponents()...)
	}

	return components
}

func (lm *LibraryManager) GetComponent(id goecs.Identifier) (cmp goecs.Component) {
	for _, library := range lm.Libraries {
		if library.GetId().Namespace == id.Namespace {
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

func (lm *LibraryManager) GetSystem(id goecs.Identifier) (sys goecs.ISystem, err error) {
	for _, library := range lm.Libraries {
		if library.GetId().Namespace == id.Namespace {
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

func (lm *LibraryManager) InstantiateSystem(id goecs.Identifier, world *goecs.IWorld) (*goecs.ISystem, error) {
	var systemLocation *goecs.ISystem

	system, err := lm.GetSystem(id)
	if err != nil {
		return nil, err
	}

	system.Init(world)

	systemLocation = &system

	return systemLocation, nil
}

func (lm *LibraryManager) InstantiateComponent(id goecs.Identifier, data interface{}) *goecs.Component {
	var componentLocation *goecs.Component

	component := lm.GetComponent(id)

	component.SetData(data)

	componentLocation = &component

	return componentLocation
}

func (lm *LibraryManager) LoadLibrary(library ILibrary) error {
	_, ok := lm.Libraries[library.GetId()]
	if !ok {
		lm.Libraries[library.GetId()] = library
	}

	return fmt.Errorf("CONFLICT: a library with same id already exist")
}

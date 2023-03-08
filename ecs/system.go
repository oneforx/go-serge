package ecs

import (
	"fmt"
)

type ISystem interface {
	GetName() string
	GetId() Identifier
	Update()
	Init(*IWorld)
	Listen(string, func(...interface{}) error) error
	Call(string, ...interface{}) error
}

type System struct {
	Id        Identifier
	Name      string
	World     *IWorld
	listening map[string]func(...interface{}) error
}

func (ss *System) Listen(id string, handler func(...interface{}) error) error {
	_, ok := ss.listening[id]
	if !ok {
		ss.listening[id] = handler
		return nil
	}

	return fmt.Errorf("the listener [%s] already exist", id)
}

func (ss *System) Call(id string, args ...interface{}) error {
	listener, ok := ss.listening[id]
	if !ok {
		return fmt.Errorf("the listener [%s] don't exist", id)
	}

	if err := listener(args...); err != nil {
		return err
	}

	return nil
}

func (ss *System) Init(world *IWorld) {
	ss.World = world
}

func (ss *System) GetName() string {
	return ss.Name
}

func (ss *System) GetId() Identifier {
	return ss.Id
}

package ecs

type ISystem interface {
	GetName() string
	GetId() Identifier
	Update()
	Init(*IWorld)
}

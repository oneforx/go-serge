package lib

import "github.com/oneforx/go-ecs"

type ILibrarySerge interface {
	ecs.ILibrary

	RegisterServerMessageHandler(ecs.Identifier, Message) error
}

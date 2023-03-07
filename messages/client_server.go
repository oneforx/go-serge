package messages

import (
	"github.com/google/uuid"
	"github.com/oneforx/go-serge/engine"
)

var (
	CS_CONNECT_TOKEN = func(message string) engine.Message {
		return engine.Message{
			MessageType: "CONNECT_TOKEN",
			Data:        message,
			Target:      engine.SERVER_TARGET,
			NetMode:     engine.NET_HYB,
		}
	}
	CS_PONG = func() engine.Message {
		return engine.Message{
			MessageType: "PONG",
			Target:      engine.SERVER_TARGET,
			NetMode:     engine.NET_TCP,
		}
	}
	CS_CREATE_CHARACTER = func() engine.Message {
		return engine.Message{
			MessageType: "CREATE_CHARACTER",
			Target:      engine.SERVER_TARGET,
			NetMode:     engine.NET_TCP,
		}
	}
	CS_PLAY_CHARACTER = func(worldId, characterId uuid.UUID) engine.Message {
		return engine.Message{
			MessageType: "PLAY_CHARACTER",
			Target:      engine.SERVER_TARGET,
			NetMode:     engine.NET_TCP,
		}
	}
	CS_DELETE_CHARACTER = func(id string) engine.Message {
		return engine.Message{
			MessageType: "DELETE_CHARACTER",
			Data:        id,
			Target:      engine.SERVER_TARGET,
			NetMode:     engine.NET_TCP,
		}
	}
	CS_CONNECT_WORLD = func() engine.Message {
		return engine.Message{
			MessageType: "CONNECT_WORLD",
			Target:      engine.SERVER_TARGET,
			NetMode:     engine.NET_TCP,
		}
	}
	CS_DISCONNECT_WORLD = func() engine.Message {
		return engine.Message{
			MessageType: "DISCONNECT_WORLD",
			Target:      engine.SERVER_TARGET,
			NetMode:     engine.NET_TCP,
		}
	}
	CS_SHOOT = func() engine.Message {
		return engine.Message{
			MessageType: "SHOOT",
			Target:      engine.SERVER_TARGET,
			NetMode:     engine.NET_TCP,
		}
	}
	CS_ORIENTATION = func(x int, y int) engine.Message {
		return engine.Message{
			MessageType: "ORIENTATION",
			Data:        map[string]interface{}{"x": x, "y": y},
			Target:      engine.SERVER_TARGET,
			NetMode:     engine.NET_UDP,
		}
	}
)

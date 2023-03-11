package messages

import (
	"github.com/google/uuid"
	"github.com/oneforx/go-serge/lib"
)

var (
	CS_CONNECT_TOKEN = func(message string) lib.Message {
		return lib.Message{
			MessageType: "CONNECT_TOKEN",
			Data:        message,
			Target:      lib.SERVER_TARGET,
			NetMode:     lib.NET_HYB,
		}
	}
	CS_PONG = func() lib.Message {
		return lib.Message{
			MessageType: "PONG",
			Target:      lib.SERVER_TARGET,
			NetMode:     lib.NET_TCP,
		}
	}
	CS_CREATE_CHARACTER = func() lib.Message {
		return lib.Message{
			MessageType: "CREATE_CHARACTER",
			Target:      lib.SERVER_TARGET,
			NetMode:     lib.NET_TCP,
		}
	}
	CS_PLAY_CHARACTER = func(worldId, characterId uuid.UUID) lib.Message {
		return lib.Message{
			MessageType: "PLAY_CHARACTER",
			Target:      lib.SERVER_TARGET,
			NetMode:     lib.NET_TCP,
		}
	}
	CS_DELETE_CHARACTER = func(id string) lib.Message {
		return lib.Message{
			MessageType: "DELETE_CHARACTER",
			Data:        id,
			Target:      lib.SERVER_TARGET,
			NetMode:     lib.NET_TCP,
		}
	}
	CS_CONNECT_WORLD = func() lib.Message {
		return lib.Message{
			MessageType: "CONNECT_WORLD",
			Target:      lib.SERVER_TARGET,
			NetMode:     lib.NET_TCP,
		}
	}
	CS_DISCONNECT_WORLD = func() lib.Message {
		return lib.Message{
			MessageType: "DISCONNECT_WORLD",
			Target:      lib.SERVER_TARGET,
			NetMode:     lib.NET_TCP,
		}
	}
	CS_SHOOT = func() lib.Message {
		return lib.Message{
			MessageType: "SHOOT",
			Target:      lib.SERVER_TARGET,
			NetMode:     lib.NET_TCP,
		}
	}
	CS_ORIENTATION = func(x int, y int) lib.Message {
		return lib.Message{
			MessageType: "ORIENTATION",
			Data:        map[string]interface{}{"x": x, "y": y},
			Target:      lib.SERVER_TARGET,
			NetMode:     lib.NET_UDP,
		}
	}
)

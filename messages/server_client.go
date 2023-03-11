package messages

import (
	"github.com/google/uuid"
	"github.com/oneforx/go-ecs"
	"github.com/oneforx/go-serge/lib"
)

var (
	SC_CONNECT_FAILED = func(data interface{}) lib.Message {
		return lib.Message{
			MessageType: "CONNECT_FAILED",
			Data:        data,
			Target:      lib.CLIENT_TARGET,
			NetMode:     lib.NET_HYB,
		}
	}
	SC_CONNECT_SUCCESS = func(data interface{}) lib.Message {
		return lib.Message{
			MessageType: "CONNECT_SUCCESS",
			Data:        data,
			Target:      lib.CLIENT_TARGET,
			NetMode:     lib.NET_HYB,
		}
	}
	SC_PING = func(latence int) lib.Message {
		return lib.Message{
			MessageType: "PING",
			Data:        latence,
			Target:      lib.CLIENT_TARGET,
			NetMode:     lib.NET_TCP,
		}
	}
	SC_CREATE_ENTITY = func(entityData ecs.EntityNoCycle) lib.Message {
		return lib.Message{
			MessageType: "CREATE_ENTITY",
			Data:        entityData,
			Target:      lib.CLIENT_TARGET,
			NetMode:     lib.NET_TCP,
		}
	}
	SC_UPDATE_ENTITY = func(entityId uuid.UUID, entityData []*ecs.Component) lib.Message {
		return lib.Message{
			MessageType: "UPDATE_ENTITY",
			Data: UpdateMessageData{
				Id:         entityId,
				Components: entityData,
			},
			Target:  lib.CLIENT_TARGET,
			NetMode: lib.NET_UDP,
		}
	}
	SC_DELETE_ENTITY = func(id string) lib.Message {
		return lib.Message{
			MessageType: "DELETE_ENTITY",
			Data:        id,
			Target:      lib.CLIENT_TARGET,
			NetMode:     lib.NET_TCP,
		}
	}
	SC_CREATE_WORLD = func(entitiesData []ecs.EntityNoCycle) lib.Message {
		return lib.Message{
			MessageType: "CREATE_WORLD",
			Data:        entitiesData,
			Target:      lib.CLIENT_TARGET,
			NetMode:     lib.NET_TCP,
		}
	}
)

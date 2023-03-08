package messages

import (
	"github.com/google/uuid"
	"github.com/oneforx/go-ecs"
	"github.com/oneforx/go-serge/engine"
)

var (
	SC_CONNECT_FAILED = func(data interface{}) engine.Message {
		return engine.Message{
			MessageType: "CONNECT_FAILED",
			Data:        data,
			Target:      engine.CLIENT_TARGET,
			NetMode:     engine.NET_HYB,
		}
	}
	SC_CONNECT_SUCCESS = func(data interface{}) engine.Message {
		return engine.Message{
			MessageType: "CONNECT_SUCCESS",
			Data:        data,
			Target:      engine.CLIENT_TARGET,
			NetMode:     engine.NET_HYB,
		}
	}
	SC_PING = func(latence int) engine.Message {
		return engine.Message{
			MessageType: "PING",
			Data:        latence,
			Target:      engine.CLIENT_TARGET,
			NetMode:     engine.NET_TCP,
		}
	}
	SC_CREATE_ENTITY = func(entityData ecs.EntityNoCycle) engine.Message {
		return engine.Message{
			MessageType: "CREATE_ENTITY",
			Data:        entityData,
			Target:      engine.CLIENT_TARGET,
			NetMode:     engine.NET_TCP,
		}
	}
	SC_UPDATE_ENTITY = func(entityId uuid.UUID, entityData []*ecs.Component) engine.Message {
		return engine.Message{
			MessageType: "UPDATE_ENTITY",
			Data: UpdateMessageData{
				Id:         entityId,
				Components: entityData,
			},
			Target:  engine.CLIENT_TARGET,
			NetMode: engine.NET_UDP,
		}
	}
	SC_DELETE_ENTITY = func(id string) engine.Message {
		return engine.Message{
			MessageType: "DELETE_ENTITY",
			Data:        id,
			Target:      engine.CLIENT_TARGET,
			NetMode:     engine.NET_TCP,
		}
	}
	SC_CREATE_WORLD = func(entitiesData []ecs.EntityNoCycle) engine.Message {
		return engine.Message{
			MessageType: "CREATE_WORLD",
			Data:        entitiesData,
			Target:      engine.CLIENT_TARGET,
			NetMode:     engine.NET_TCP,
		}
	}
)

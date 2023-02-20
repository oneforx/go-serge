package systems

import (
	"log"
	"sync"

	"github.com/bawdeveloppement/squirrelgameserver/ecs"
	"github.com/bawdeveloppement/squirrelgameserver/engine"
	"github.com/bawdeveloppement/squirrelgameserver/messages"
	"github.com/google/uuid"
)

type NetUpdateSystem struct {
	Id           uuid.UUID
	Name         string
	World        *ecs.IWorld
	ServerEngine *engine.ServerEngine
}

func (ss *NetUpdateSystem) GetName() string {
	return ss.Name
}

func (ss *NetUpdateSystem) GetId() uuid.UUID {
	return ss.Id
}

func (ss *NetUpdateSystem) Update(mutex *sync.Mutex) {

	if ss.ServerEngine.UdpInit {
		worldLocalised := *ss.World

		for _, entity := range worldLocalised.GetEntities() {
			entityLocalised := *entity

			bytes, err := engine.MessageToBytes(messages.SC_UPDATE_ENTITY(entityLocalised.GetId(), entityLocalised.GetComponents()))
			if err != nil {
				log.Println("Could not parse SC_UPDATE_ENTITY to bytes")
				continue
			}

			for _, u := range ss.ServerEngine.UdpConnexions {
				log.Println("dzadaz")
				_, err = ss.ServerEngine.UdpListener.WriteToUDP(bytes, u)
				if err != nil {
					log.Println("Could not SC_UPDATE_ENTITY")
					continue
				}
			}
		}
	}
}

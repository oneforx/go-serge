package systems

import (
	"fmt"
	"log"
	"sync"

	"github.com/oneforx/go-serge/ecs"
	"github.com/oneforx/go-serge/engine"
	"github.com/oneforx/go-serge/messages"
)

type NetUpdateSystem struct {
	ecs.System
	Id           ecs.Identifier
	Name         string
	World        *ecs.IWorld
	ServerEngine *engine.ServerEngine
	listening    map[string]func(...interface{}) error
}

func (ss *NetUpdateSystem) Listen(id string, handler func(...interface{}) error) error {
	_, ok := ss.listening[id]
	if !ok {
		ss.listening[id] = handler
		return nil
	}

	return fmt.Errorf("the listener [%s] already exist", id)
}

func (ss *NetUpdateSystem) Call(id string, args ...interface{}) error {
	listener, ok := ss.listening[id]
	if !ok {
		return fmt.Errorf("the listener [%s] don't exist", id)
	}

	if err := listener(args...); err != nil {
		return err
	}

	return nil
}

func (ss *NetUpdateSystem) Init(world *ecs.IWorld) {
	ss.World = world
}

func (ss *NetUpdateSystem) GetName() string {
	return ss.Name
}

func (ss *NetUpdateSystem) GetId() ecs.Identifier {
	return ss.Id
}

func (ss *NetUpdateSystem) Update() {

	var sync_mutex sync.Mutex
	sync_mutex.Lock()
	init := ss.ServerEngine.UdpInit
	sync_mutex.Unlock()
	if init {
		worldLocalised := *ss.World

		for _, entity := range worldLocalised.GetEntities() {
			entityLocalised := *entity

			var components []*ecs.Component = []*ecs.Component{}

			for _, cmp := range entityLocalised.GetComponents() {
				cmpLocalised := *cmp
				components = append(components, cmpLocalised.GetStructure())
			}

			bytes, err := engine.MessageToBytes(
				messages.SC_UPDATE_ENTITY(entityLocalised.GetId(), components),
			)
			if err != nil {
				log.Println("Could not parse SC_UPDATE_ENTITY to bytes")
				continue
			}

			sync_mutex.Lock()
			udpConnexions := ss.ServerEngine.UdpConnexions
			sync_mutex.Unlock()
			for _, u := range udpConnexions {
				sync_mutex.Lock()
				_, err = ss.ServerEngine.UdpListener.WriteToUDP(bytes, u)
				if err != nil {
					defer sync_mutex.Unlock()
					log.Println("Could not SC_UPDATE_ENTITY")
					continue
				}
				sync_mutex.Unlock()
			}
		}
	}
}

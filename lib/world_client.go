package lib

import "github.com/oneforx/go-ecs"

type WorldClient struct {
	World *ecs.IWorld

	States map[string]interface{}
}

func (wc *WorldClient) TSend(message Message) {

}

func (wc *WorldClient) USend(message Message) {

}

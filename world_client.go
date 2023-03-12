package goserge

import goecs "github.com/oneforx/go-ecs"

type WorldClient struct {
	World *goecs.IWorld

	States map[string]interface{}
}

func (wc *WorldClient) TSend(message Message) {

}

func (wc *WorldClient) USend(message Message) {

}

package main

import (
	"github.com/oneforx/go-serge/ecs"
	"github.com/oneforx/go-serge/library"
)

func Load() library.ILibrary {
	var myLibrary = &library.Library{
		Id: ecs.Identifier{
			Namespace: "oneforx",
			Path:      "base",
		},
	}

	var bulletSystem ecs.ISystem = &BulletSystem{
		Id: ecs.Identifier{
			Namespace: "oneforx",
			Path:      "bullet",
		},
	}

	myLibrary.AddSystem(bulletSystem)

	myLibrary.AddComponent(ecs.Component{
		Id: ecs.Identifier{
			Namespace: "oneforx",
			Path:      "position",
		},
		Data: map[string]interface{}{"x": 0, "y": 0, "origin_x": 0, "origin_y": 0, "vel_x": 0, "vel_y": 0},
	})

	myLibrary.AddComponent(ecs.Component{
		Id: ecs.Identifier{
			Namespace: "oneforx",
			Path:      "dimension",
		},
		Data: map[string]interface{}{"width": 0, "height": 0},
	})

	myLibrary.AddComponent(ecs.Component{
		Id: ecs.Identifier{
			Namespace: "oneforx",
			Path:      "solid",
		},
		Data: false,
	})

	myLibrary.AddComponent(ecs.Component{
		Id: ecs.Identifier{
			Namespace: "oneforx",
			Path:      "speed",
		},
		Data: 0,
	})

	myLibrary.AddComponent(ecs.Component{
		Id: ecs.Identifier{
			Namespace: "oneforx",
			Path:      "force",
		},
		Data: 0,
	})

	myLibrary.AddComponent(ecs.Component{
		Id: ecs.Identifier{
			Namespace: "oneforx",
			Path:      "orientation",
		},
		Data: 0,
	})

	myLibrary.AddComponent(ecs.Component{
		Id: ecs.Identifier{
			Namespace: "oneforx",
			Path:      "distance",
		},
		Data: 0,
	})
	return myLibrary
}

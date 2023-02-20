package messages

import (
	"github.com/bawdeveloppement/squirrelgameserver/ecs"
	"github.com/google/uuid"
)

type UpdateMessageData struct {
	Id         uuid.UUID        `json:"id"`
	Components []*ecs.Component `json:"components"`
}

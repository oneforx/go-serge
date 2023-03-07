package messages

import (
	"github.com/google/uuid"
	"github.com/oneforx/go-serge/ecs"
)

type UpdateMessageData struct {
	Id         uuid.UUID        `json:"id"`
	Components []*ecs.Component `json:"components"`
}

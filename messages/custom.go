package messages

import (
	"github.com/google/uuid"
	goecs "github.com/oneforx/go-ecs"
)

type UpdateMessageData struct {
	Id         uuid.UUID          `json:"id"`
	Components []*goecs.Component `json:"components"`
}

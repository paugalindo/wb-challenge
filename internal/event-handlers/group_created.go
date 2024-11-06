package eventhandlers

import (
	"context"
	"encoding/json"
	"log"
	"wb-challenge/bus"
	"wb-challenge/internal"
	"wb-challenge/internal/commands"
)

type GroupCreated struct {
	cmdBus *bus.CommandBus

	logger *log.Logger
}

func NewGroupCreated(cmdBus *bus.CommandBus, logger *log.Logger) GroupCreated {
	return GroupCreated{
		cmdBus: cmdBus,
		logger: logger,
	}
}

func (h *GroupCreated) Handle(data []byte) {
	e := internal.GroupCreatedEvent{}
	if err := json.Unmarshal(data, &e); err != nil {
		h.logger.Printf("error handling GroupCreatedEvent: %e", err)
	}

	cmd := commands.AssignVehicleToGroupCmd{
		GroupID: e.ID,
	}
	if err := h.cmdBus.Dispatch(context.Background(), cmd); err != nil {
		h.logger.Printf("error assigning vehicle on GroupCreatedEvent: %e", err)
	}
}

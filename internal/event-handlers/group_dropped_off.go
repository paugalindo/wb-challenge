package eventhandlers

import (
	"context"
	"encoding/json"
	"log"
	"wb-challenge/bus"
	"wb-challenge/internal"
	"wb-challenge/internal/commands"
)

type GroupDroppedOff struct {
	cmdBus *bus.CommandBus

	logger *log.Logger
}

func NewGroupDroppedOff(cmdBus *bus.CommandBus, logger *log.Logger) GroupDroppedOff {
	return GroupDroppedOff{
		cmdBus: cmdBus,
		logger: logger,
	}
}

func (h *GroupDroppedOff) Handle(data []byte) {
	e := internal.GroupDroppedOffEvent{}
	if err := json.Unmarshal(data, &e); err != nil {
		h.logger.Printf("error handling GroupDroppedOffEvent: %e", err)
	}

	if err := h.cmdBus.Dispatch(context.Background(), commands.ReleaseVehicleCmd{GroupID: e.ID}); err != nil {
		h.logger.Printf("error assigning vehicle on GroupDroppedOffEvent: %e", err)
	}
}

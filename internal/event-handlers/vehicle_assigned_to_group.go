package eventhandlers

import (
	"context"
	"encoding/json"
	"log"
	"wb-challenge/bus"
	"wb-challenge/internal"
	"wb-challenge/internal/commands"
)

type VehicleAssignedToGroup struct {
	cmdBus *bus.CommandBus

	logger *log.Logger
}

func NewVehicleAssignedToGroup(cmdBus *bus.CommandBus, logger *log.Logger) VehicleAssignedToGroup {
	return VehicleAssignedToGroup{
		cmdBus: cmdBus,
		logger: logger,
	}
}

func (h *VehicleAssignedToGroup) Handle(data []byte) {
	e := internal.VehicleAssignedToGroupEvent{}
	if err := json.Unmarshal(data, &e); err != nil {
		h.logger.Printf("error handling VehicleAssignedToGroupEvent: %e", err)
	}

	cmd := commands.OccupyVehicleCmd{
		ID:      e.VehicleID,
		GroupID: e.ID,
	}
	if err := h.cmdBus.Dispatch(context.Background(), cmd); err != nil {
		h.logger.Printf("error assigning vehicle on VehicleAssignedToGroupEvent: %e", err)
	}
}

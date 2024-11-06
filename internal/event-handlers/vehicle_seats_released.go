package eventhandlers

import (
	"context"
	"encoding/json"
	"log"
	"wb-challenge/bus"
	"wb-challenge/internal"
	"wb-challenge/internal/commands"
)

type VehicleSeatsReleased struct {
	cmdBus *bus.CommandBus

	logger *log.Logger
}

func NewVehicleSeatsReleased(cmdBus *bus.CommandBus, logger *log.Logger) VehicleSeatsReleased {
	return VehicleSeatsReleased{
		cmdBus: cmdBus,
		logger: logger,
	}
}

func (h *VehicleSeatsReleased) Handle(data []byte) {
	e := internal.VehicleSeatsReleasedEvent{}
	if err := json.Unmarshal(data, &e); err != nil {
		h.logger.Printf("error handling VehicleSeatsReleasedEvent: %e", err)
	}

	if err := h.cmdBus.Dispatch(context.Background(), commands.AssignVehiclesCmd{}); err != nil {
		h.logger.Printf("error assigning vehicle on VehicleSeatsReleasedEvent: %e", err)
	}
}

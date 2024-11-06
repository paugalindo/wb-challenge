package commands

import (
	"context"
	"errors"
	"sync"
	"wb-challenge/internal"
)

const OccupyVehicleType = "commands.OccupyVehicleCmd"

type OccupyVehicleCmd struct {
	ID      int
	GroupID int
}

type OccupyVehicleHandler struct {
	mutex             *sync.Mutex
	groupView         internal.GroupView
	vehicleRepository internal.VehicleRepository
	publisher         internal.EventsPublisher
}

func NewOccupyVehicleHandler(mutex *sync.Mutex, groupView internal.GroupView,
	vehicleRepository internal.VehicleRepository, publisher internal.EventsPublisher,
) OccupyVehicleHandler {
	return OccupyVehicleHandler{
		mutex:             mutex,
		groupView:         groupView,
		vehicleRepository: vehicleRepository,
		publisher:         publisher,
	}
}

func (h *OccupyVehicleHandler) Handle(ctx context.Context, c any) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	cmd, ok := c.(OccupyVehicleCmd)
	if !ok {
		return errors.New("invalid command type")
	}

	v, err := h.vehicleRepository.Get(cmd.ID)
	if err != nil {
		return nil
	}

	g, err := h.groupView.Get(cmd.GroupID)
	if err != nil {
		return err
	}

	if err := v.OccupySeats(g.TotalPeople()); err != nil {
		return err
	}

	if err := h.vehicleRepository.Save(v); err != nil {
		return err
	}

	return h.publisher.Publish(v.Events()...)
}

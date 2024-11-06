package commands

import (
	"context"
	"errors"
	"sync"
	"wb-challenge/internal"
)

const ReleaseVehicleType = "commands.ReleaseVehicleCmd"

type ReleaseVehicleCmd struct {
	GroupID int
}

type ReleaseVehicleHandler struct {
	mutex             *sync.Mutex
	groupView         internal.GroupView
	vehicleRepository internal.VehicleRepository
	publisher         internal.EventsPublisher
}

func NewReleaseVehicleHandler(mutex *sync.Mutex, groupView internal.GroupView,
	vehicleRepository internal.VehicleRepository, publisher internal.EventsPublisher,
) ReleaseVehicleHandler {
	return ReleaseVehicleHandler{
		mutex:             mutex,
		groupView:         groupView,
		vehicleRepository: vehicleRepository,
		publisher:         publisher,
	}
}

func (h *ReleaseVehicleHandler) Handle(ctx context.Context, c any) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	cmd, ok := c.(ReleaseVehicleCmd)
	if !ok {
		return errors.New("invalid command type")
	}

	g, err := h.groupView.Get(cmd.GroupID)
	if err != nil {
		if errors.Is(err, internal.ErrGroupNotFound) {
			return ErrNotFound
		}
		return err
	}

	v, err := h.vehicleRepository.Get(g.VehicleAssigned())
	if err != nil {
		return err
	}

	if err := v.ReleaseSeats(g.TotalPeople()); err != nil {
		return err
	}

	if err = h.vehicleRepository.Save(v); err != nil {
		return err
	}

	return h.publisher.Publish(v.Events()...)
}

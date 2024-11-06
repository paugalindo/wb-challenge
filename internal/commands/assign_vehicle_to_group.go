package commands

import (
	"context"
	"errors"
	"sync"
	"wb-challenge/internal"
)

const AssignVehicleToGroupType = "commands.AssignVehicleToGroupCmd"

type AssignVehicleToGroupCmd struct {
	GroupID int
}

type AssignVehicleToGroupHandler struct {
	mutex           *sync.Mutex
	groupRepository internal.GroupRepository
	vehicleView     internal.VehicleView
	publisher       internal.EventsPublisher
}

func NewAssignVehicleToGroupHandler(mutex *sync.Mutex, groupRepository internal.GroupRepository,
	vehicleView internal.VehicleView, publisher internal.EventsPublisher,
) AssignVehicleToGroupHandler {
	return AssignVehicleToGroupHandler{
		mutex:           mutex,
		groupRepository: groupRepository,
		vehicleView:     vehicleView,
		publisher:       publisher,
	}
}

func (h *AssignVehicleToGroupHandler) Handle(ctx context.Context, c any) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	cmd, ok := c.(AssignVehicleToGroupCmd)
	if !ok {
		return errors.New("invalid command type")
	}

	g, err := h.groupRepository.Get(cmd.GroupID)
	if err != nil {
		return err
	}

	if g.VehicleAssigned() != 0 || g.IsDroppedOff() {
		return nil
	}

	v, err := h.vehicleView.GetWithEmptySeats(g.TotalPeople())
	if err != nil {
		return nil
	}

	if err := g.AssignVehicle(v.ID()); err != nil {
		return err
	}

	if err := h.groupRepository.Save(g); err != nil {
		return err
	}

	return h.publisher.Publish(g.Events()...)
}

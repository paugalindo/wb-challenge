package commands

import (
	"context"
	"errors"
	"sync"
	"wb-challenge/internal"
)

const AssignVehiclesType = "commands.AssignVehiclesCmd"

type AssignVehiclesCmd struct {
	GroupID int
}

type AssignVehiclesHandler struct {
	mutex           *sync.Mutex
	groupRepository internal.GroupRepository
	vehicleView     internal.VehicleView
	publisher       internal.EventsPublisher
}

func NewAssignVehiclesHandler(mutex *sync.Mutex, groupRepository internal.GroupRepository,
	vehicleView internal.VehicleView, publisher internal.EventsPublisher,
) AssignVehiclesHandler {
	return AssignVehiclesHandler{
		mutex:           mutex,
		groupRepository: groupRepository,
		vehicleView:     vehicleView,
		publisher:       publisher,
	}
}

func (h *AssignVehiclesHandler) Handle(ctx context.Context, c any) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := c.(AssignVehiclesCmd); !ok {
		return errors.New("invalid command type")
	}

	groups, err := h.groupRepository.GetUnassignedOrderedByCreatedAt()
	if err != nil {
		return err
	}

	for _, g := range groups {
		if g.VehicleAssigned() != 0 || g.IsDroppedOff() {
			continue
		}

		v, err := h.vehicleView.GetWithEmptySeats(g.TotalPeople())
		if err != nil {
			continue
		}

		if err := g.AssignVehicle(v.ID()); err != nil {
			return err
		}

		if err := h.groupRepository.Save(g); err != nil {
			return err
		}

		if err := h.publisher.Publish(g.Events()...); err != nil {
			return err
		}
	}

	return nil
}

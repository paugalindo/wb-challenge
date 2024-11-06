package commands

import (
	"context"
	"errors"
	"sync"
	"wb-challenge/internal"
)

const LoadVehiclesType = "commands.LoadVehiclesCmd"

type LoadVehiclesCmd struct {
	Vehicles []Vehicle
}

type Vehicle struct {
	ID    int
	Seats int
}

type LoadVehiclesHandler struct {
	mutex             *sync.Mutex
	vehicleRepository internal.VehicleRepository
	groupRepository   internal.GroupRepository
	publisher         internal.EventsPublisher
}

func NewLoadVehiclesHandler(mutex *sync.Mutex, vehicleRepository internal.VehicleRepository, groupRepository internal.GroupRepository, publisher internal.EventsPublisher) LoadVehiclesHandler {
	return LoadVehiclesHandler{
		mutex:             mutex,
		vehicleRepository: vehicleRepository,
		groupRepository:   groupRepository,
		publisher:         publisher,
	}
}

func (h *LoadVehiclesHandler) Handle(ctx context.Context, c any) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	cmd, ok := c.(LoadVehiclesCmd)
	if !ok {
		return errors.New("invalid command type")
	}

	if err := h.vehicleRepository.RemoveAllVehicles(); err != nil {
		return err
	}

	if err := h.groupRepository.RemoveAllGroups(); err != nil {
		return err
	}

	for _, v := range cmd.Vehicles {
		v, err := internal.NewVehicle(v.ID, v.Seats)
		if err != nil {
			return err
		}

		if err = h.vehicleRepository.Save(v); err != nil {
			return err
		}

		if err = h.publisher.Publish(v.Events()...); err != nil {
			return err
		}
	}
	return nil
}

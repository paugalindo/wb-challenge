package commands_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"wb-challenge/internal"
	"wb-challenge/internal/commands"
)

func TestAssignVehiclesHandler_Handle(t *testing.T) {
	groupID1, vehicleID1, groupPeople1, vehicleSeats1 := 1, 1, 1, 1
	groupID2, vehicleID2, groupPeople2, vehicleSeats2 := 2, 2, 2, 2
	mockGroupRepo := internal.NewMockGroupRepository()
	mockGroupRepo.GetUnassignedOrderedByCreatedAtFunc = func() ([]internal.Group, error) {
		return []internal.Group{
			internal.HydrateGroup(groupID1, groupPeople1, 0, false),
			internal.HydrateGroup(groupID2, groupPeople2, 0, false),
		}, nil
	}
	mockGroupRepo.SaveFunc = func(internal.Group) error {
		return nil
	}
	mockVehicleRepo := internal.NewMockVehicleRepository()
	mockVehicleRepo.GetWithEmptySeatsFunc = func(seats int) (internal.Vehicle, error) {
		switch seats {
		case groupPeople1:
			return internal.HydrateVehicle(vehicleID1, vehicleSeats1, 0), nil
		case groupPeople2:
			return internal.HydrateVehicle(vehicleID2, vehicleSeats2, 0), nil
		}
		return internal.Vehicle{}, errors.New("error")
	}
	mockPublisher := internal.NewMockEventsPublisher()
	mockPublisher.PublishFunc = func(...internal.Event) error {
		return nil
	}

	handler := commands.NewAssignVehiclesHandler(
		&sync.Mutex{},
		mockGroupRepo,
		mockVehicleRepo,
		mockPublisher,
	)

	cmd := commands.AssignVehiclesCmd{GroupID: groupID1}

	err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mockGroupRepo.GetUnassignedOrderedByCreatedAtCalls) != 1 {
		t.Errorf("Expected 1 call to GroupRepository.Get, got %d", len(mockGroupRepo.GetUnassignedOrderedByCreatedAtCalls))
	}

	if len(mockVehicleRepo.GetWithEmptySeatsCalls) != 2 {
		t.Errorf("Expected 2 call to VehicleRepository.GetWithEmptySeats, got %d", len(mockVehicleRepo.GetWithEmptySeatsCalls))
	}
	if mockVehicleRepo.GetWithEmptySeatsCalls[0].Seats != groupPeople1 {
		t.Errorf("Expected empty seats %d, got %d", groupPeople1, mockVehicleRepo.GetWithEmptySeatsCalls[0].Seats)
	}

	if len(mockGroupRepo.SaveCalls) != 2 {
		t.Errorf("Expected 2 call to GroupRepository.Save, got %d", len(mockGroupRepo.SaveCalls))
	}
	if mockGroupRepo.SaveCalls[0].Group.VehicleAssigned() != vehicleID1 {
		t.Errorf("Expected vehicle ID %d, got %d", vehicleID1, mockGroupRepo.SaveCalls[0].Group.VehicleAssigned())
	}
	if mockGroupRepo.SaveCalls[1].Group.VehicleAssigned() != vehicleID2 {
		t.Errorf("Expected vehicle ID %d, got %d", vehicleID2, mockGroupRepo.SaveCalls[0].Group.VehicleAssigned())
	}

	if len(mockPublisher.PublishCalls) != 2 {
		t.Errorf("Expected 1 call to EventPublisher.Publish, got %d", len(mockPublisher.PublishCalls))
	}
	if len(mockPublisher.PublishCalls[0].Events) != 1 {
		t.Errorf("Expected 1 events published, got %d", len(mockPublisher.PublishCalls[0].Events))
	}
	if mockPublisher.PublishCalls[0].Events[0].Type() != internal.VehicleAssignedToGroupEventType {
		t.Errorf("Expected %s event, got %s", internal.VehicleAssignedToGroupEventType, mockPublisher.PublishCalls[0].Events[0].Type())
	}
	if len(mockPublisher.PublishCalls[1].Events) != 1 {
		t.Errorf("Expected 1 events published, got %d", len(mockPublisher.PublishCalls[0].Events))
	}
	if mockPublisher.PublishCalls[1].Events[0].Type() != internal.VehicleAssignedToGroupEventType {
		t.Errorf("Expected %s event, got %s", internal.VehicleAssignedToGroupEventType, mockPublisher.PublishCalls[1].Events[0].Type())
	}
}

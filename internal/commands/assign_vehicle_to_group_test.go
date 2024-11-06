package commands_test

import (
	"context"
	"sync"
	"testing"
	"wb-challenge/internal"
	"wb-challenge/internal/commands"
)

func TestAssignVehicleToGroupHandler_Handle(t *testing.T) {
	groupID, vehicleID, groupPeople, vehicleSeats := 1, 1, 1, 1
	mockGroupRepo := internal.NewMockGroupRepository()
	mockGroupRepo.GetFunc = func(id int) (internal.Group, error) {
		return internal.HydrateGroup(groupID, groupPeople, 0, false), nil
	}
	mockGroupRepo.SaveFunc = func(internal.Group) error {
		return nil
	}
	mockVehicleRepo := internal.NewMockVehicleRepository()
	mockVehicleRepo.GetWithEmptySeatsFunc = func(seats int) (internal.Vehicle, error) {
		return internal.HydrateVehicle(vehicleID, vehicleSeats, 0), nil
	}
	mockPublisher := internal.NewMockEventsPublisher()
	mockPublisher.PublishFunc = func(...internal.Event) error {
		return nil
	}

	handler := commands.NewAssignVehicleToGroupHandler(
		&sync.Mutex{},
		mockGroupRepo,
		mockVehicleRepo,
		mockPublisher,
	)

	cmd := commands.AssignVehicleToGroupCmd{GroupID: groupID}

	err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mockGroupRepo.GetCalls) != 1 {
		t.Errorf("Expected 1 call to GroupRepository.Get, got %d", len(mockGroupRepo.GetCalls))
	}
	if mockGroupRepo.GetCalls[0].ID != groupID {
		t.Errorf("Expected group ID %d, got %d", groupID, mockGroupRepo.GetCalls[0].ID)
	}

	if len(mockVehicleRepo.GetWithEmptySeatsCalls) != 1 {
		t.Errorf("Expected 1 call to VehicleRepository.GetWithEmptySeats, got %d", len(mockVehicleRepo.GetWithEmptySeatsCalls))
	}
	if mockVehicleRepo.GetWithEmptySeatsCalls[0].Seats != groupPeople {
		t.Errorf("Expected empty seats %d, got %d", groupPeople, mockVehicleRepo.GetWithEmptySeatsCalls[0].Seats)
	}

	if len(mockGroupRepo.SaveCalls) != 1 {
		t.Errorf("Expected 1 call to GroupRepository.Save, got %d", len(mockGroupRepo.SaveCalls))
	}
	if mockGroupRepo.SaveCalls[0].Group.VehicleAssigned() != vehicleID {
		t.Errorf("Expected vehicle ID %d, got %d", vehicleID, mockGroupRepo.SaveCalls[0].Group.VehicleAssigned())
	}

	if len(mockPublisher.PublishCalls) != 1 {
		t.Errorf("Expected 1 call to EventPublisher.Publish, got %d", len(mockPublisher.PublishCalls))
	}
	if len(mockPublisher.PublishCalls[0].Events) != 1 {
		t.Errorf("Expected 1 events published, got %d", len(mockPublisher.PublishCalls[0].Events))
	}
	if mockPublisher.PublishCalls[0].Events[0].Type() != internal.VehicleAssignedToGroupEventType {
		t.Errorf("Expected %s event, got %s", internal.VehicleAssignedToGroupEventType, mockPublisher.PublishCalls[0].Events[0].Type())
	}
}

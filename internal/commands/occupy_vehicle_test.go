package commands_test

import (
	"context"
	"sync"
	"testing"
	"wb-challenge/internal"
	"wb-challenge/internal/commands"
)

func TestOccupyVehicleHandler_Handle(t *testing.T) {
	groupID, vehicleID, groupPeople, vehicleSeats := 1, 1, 1, 1
	mockGroupRepo := internal.NewMockGroupRepository()
	mockGroupRepo.GetFunc = func(id int) (internal.Group, error) {
		return internal.HydrateGroup(groupID, groupPeople, vehicleID, false), nil
	}
	mockVehicleRepo := internal.NewMockVehicleRepository()
	mockVehicleRepo.GetFunc = func(seats int) (internal.Vehicle, error) {
		return internal.HydrateVehicle(vehicleID, vehicleSeats, 0), nil
	}
	mockPublisher := internal.NewMockEventsPublisher()
	mockPublisher.PublishFunc = func(...internal.Event) error {
		return nil
	}

	handler := commands.NewOccupyVehicleHandler(
		&sync.Mutex{},
		mockGroupRepo,
		mockVehicleRepo,
		mockPublisher,
	)

	cmd := commands.OccupyVehicleCmd{ID: vehicleID, GroupID: groupID}

	err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mockVehicleRepo.GetCalls) != 1 {
		t.Errorf("Expected 1 call to VehicleRepository.Get, got %d", len(mockVehicleRepo.GetCalls))
	}
	if mockVehicleRepo.GetCalls[0].ID != vehicleID {
		t.Errorf("Expected vehicle id %d, got %d", vehicleID, mockVehicleRepo.GetCalls[0].ID)
	}

	if len(mockGroupRepo.GetCalls) != 1 {
		t.Errorf("Expected 1 call to GroupRepository.Get, got %d", len(mockGroupRepo.GetCalls))
	}
	if mockGroupRepo.GetCalls[0].ID != groupID {
		t.Errorf("Expected group ID %d, got %d", groupID, mockGroupRepo.GetCalls[0].ID)
	}

	if len(mockVehicleRepo.SaveCalls) != 1 {
		t.Errorf("Expected 1 call to VehicleRepository.Save, got %d", len(mockGroupRepo.SaveCalls))
	}
	if mockVehicleRepo.SaveCalls[0].Vehicle.OccupiedSeats() != groupPeople {
		t.Errorf("Expected occupied seats %d, got %d", groupPeople, mockVehicleRepo.SaveCalls[0].Vehicle.OccupiedSeats())
	}

	if len(mockPublisher.PublishCalls) != 1 {
		t.Errorf("Expected 1 call to EventPublisher.Publish, got %d", len(mockPublisher.PublishCalls))
	}
	if len(mockPublisher.PublishCalls[0].Events) != 1 {
		t.Errorf("Expected 1 events published, got %d", len(mockPublisher.PublishCalls[0].Events))
	}
	if mockPublisher.PublishCalls[0].Events[0].Type() != internal.VehicleSeatsOccupiedEventType {
		t.Errorf("Expected %s event, got %s", internal.VehicleSeatsOccupiedEventType, mockPublisher.PublishCalls[0].Events[0].Type())
	}
}

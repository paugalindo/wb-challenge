package commands_test

import (
	"context"
	"sync"
	"testing"
	"wb-challenge/internal"
	"wb-challenge/internal/commands"
)

func TestLoadVehiclesHandler_Handle(t *testing.T) {
	mockGroupRepo := internal.NewMockGroupRepository()
	mockGroupRepo.RemoveAllFunc = func() error {
		return nil
	}
	mockVehicleRepo := internal.NewMockVehicleRepository()
	mockVehicleRepo.RemoveAllFunc = func() error {
		return nil
	}
	mockVehicleRepo.SaveFunc = func(internal.Vehicle) error {
		return nil
	}
	mockPublisher := internal.NewMockEventsPublisher()
	mockPublisher.PublishFunc = func(...internal.Event) error {
		return nil
	}

	handler := commands.NewLoadVehiclesHandler(
		&sync.Mutex{},
		mockVehicleRepo,
		mockGroupRepo,
		mockPublisher,
	)

	vehicleID1, vehicleID2, seats1, seats2 := 1, 2, 4, 6
	cmd := commands.LoadVehiclesCmd{
		Vehicles: []commands.Vehicle{
			{ID: vehicleID1, Seats: seats1},
			{ID: vehicleID2, Seats: seats2},
		},
	}

	err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mockVehicleRepo.RemoveAllCalls) != 1 {
		t.Errorf("Expected 1 call to VehicleRepository.RemoveAll, got %d", len(mockVehicleRepo.RemoveAllCalls))
	}
	if len(mockGroupRepo.RemoveAllCalls) != 1 {
		t.Errorf("Expected 1 call to GroupRepository.RemoveAll, got %d", len(mockGroupRepo.RemoveAllCalls))
	}

	if len(mockVehicleRepo.SaveCalls) != 2 {
		t.Errorf("Expected 2 call to VehicleRepository.Save, got %d", len(mockGroupRepo.SaveCalls))
	}
	if mockVehicleRepo.SaveCalls[0].Vehicle.ID() != vehicleID1 {
		t.Errorf("Expected vehicleID %d, got %d", vehicleID1, mockVehicleRepo.SaveCalls[0].Vehicle.ID())
	}
	if mockVehicleRepo.SaveCalls[0].Vehicle.AvailableSeats() != seats1 {
		t.Errorf("Expected available seats %d, got %d", seats1, mockVehicleRepo.SaveCalls[0].Vehicle.AvailableSeats())
	}
	if mockVehicleRepo.SaveCalls[1].Vehicle.ID() != vehicleID2 {
		t.Errorf("Expected vehicleID %d, got %d", vehicleID2, mockVehicleRepo.SaveCalls[1].Vehicle.ID())
	}
	if mockVehicleRepo.SaveCalls[1].Vehicle.AvailableSeats() != seats2 {
		t.Errorf("Expected available seats %d, got %d", seats1, mockVehicleRepo.SaveCalls[1].Vehicle.AvailableSeats())
	}

	if len(mockPublisher.PublishCalls) != 2 {
		t.Errorf("Expected 1 call to EventPublisher.Publish, got %d", len(mockPublisher.PublishCalls))
	}
	if len(mockPublisher.PublishCalls[0].Events) != 1 {
		t.Errorf("Expected 1 events published, got %d", len(mockPublisher.PublishCalls[0].Events))
	}
	if mockPublisher.PublishCalls[0].Events[0].Type() != internal.VehicleCreatedEventType {
		t.Errorf("Expected %s event, got %s", internal.VehicleCreatedEventType, mockPublisher.PublishCalls[0].Events[0].Type())
	}
	if len(mockPublisher.PublishCalls[1].Events) != 1 {
		t.Errorf("Expected 2 events published, got %d", len(mockPublisher.PublishCalls[0].Events))
	}
	if mockPublisher.PublishCalls[1].Events[0].Type() != internal.VehicleCreatedEventType {
		t.Errorf("Expected %s event, got %s", internal.VehicleCreatedEventType, mockPublisher.PublishCalls[1].Events[0].Type())
	}
}

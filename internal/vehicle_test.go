package internal_test

import (
	"testing"
	"wb-challenge/internal"
)

func TestNewVehicle(t *testing.T) {
	vehicleID, seats := 1, 4
	v, err := internal.NewVehicle(vehicleID, seats)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if v.ID() != vehicleID {
		t.Errorf("Expected vehicle ID %d, got %d", vehicleID, v.ID())
	}
	if v.AvailableSeats() != seats {
		t.Errorf("Expected %d seats, got %d", seats, v.AvailableSeats())
	}
	if len(v.Events()) != 1 {
		t.Errorf("Expected 1 event, got %d", len(v.Events()))
	}
	if v.Events()[0].Type() != internal.VehicleCreatedEventType {
		t.Errorf("Expected %s event, got %s", internal.VehicleCreatedEventType, v.Events()[0].Type())
	}

	_, err = internal.NewVehicle(0, seats)
	if err == nil || err.Error() != "wrong ID" {
		t.Fatalf("Expected error 'wrong ID', got %v", err)
	}

	_, err = internal.NewVehicle(vehicleID, 3)
	if err == nil || err.Error() != "wrong number of seats" {
		t.Fatalf("Expected error 'wrong number of seats', got %v", err)
	}
}

func TestHydrateVehicle(t *testing.T) {
	vehicleID, totalSeats, occupiedSeats := 2, 5, 2
	v := internal.HydrateVehicle(vehicleID, totalSeats, occupiedSeats)
	if v.ID() != vehicleID {
		t.Errorf("Expected vehicle ID %d, got %d", vehicleID, v.ID())
	}
	if v.OccupiedSeats() != occupiedSeats {
		t.Errorf("Expected %d occupied seats, got %d", occupiedSeats, v.OccupiedSeats())
	}
	availableSeats := totalSeats - occupiedSeats
	if v.AvailableSeats() != availableSeats {
		t.Errorf("Expected %d available seats, got %d", availableSeats, v.AvailableSeats())
	}
	if len(v.Events()) != 0 {
		t.Errorf("Expected empty events, got %d", len(v.Events()))
	}
}

func TestVehicle_OccupySeats(t *testing.T) {
	vehicleID, seats := 3, 5
	v, _ := internal.NewVehicle(vehicleID, seats)

	seatsToOccupy := 3
	err := v.OccupySeats(seatsToOccupy)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if v.OccupiedSeats() != seatsToOccupy {
		t.Errorf("Expected %d occupied seats, got %d", seatsToOccupy, v.OccupiedSeats())
	}
	availableSeats := seats - seatsToOccupy
	if v.AvailableSeats() != availableSeats {
		t.Errorf("Expected %d available seats, got %d", availableSeats, v.AvailableSeats())
	}

	err = v.OccupySeats(seatsToOccupy)
	if err == nil || err.Error() != "not enought available seats" {
		t.Fatalf("Expected error 'not enought available seats', got %v", err)
	}

	if len(v.Events()) != 2 {
		t.Errorf("Expected 2 events, got %d", len(v.Events()))
	}
	if v.Events()[1].Type() != internal.VehicleSeatsOccupiedEventType {
		t.Errorf("Expected %s event, got %s", internal.VehicleSeatsOccupiedEventType, v.Events()[1].Type())
	}
}

func TestVehicle_ReleaseSeats(t *testing.T) {
	vehicleID, totalSeats, occupiedSeats := 4, 5, 4
	v := internal.HydrateVehicle(vehicleID, totalSeats, occupiedSeats)

	seatsToRelease := 2
	err := v.ReleaseSeats(seatsToRelease)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	occupiedSeatsAfterRelease := occupiedSeats - seatsToRelease
	if v.OccupiedSeats() != occupiedSeatsAfterRelease {
		t.Errorf("Expected %d occupied seats, got %d", occupiedSeatsAfterRelease, v.OccupiedSeats())
	}
	availableSeats := totalSeats - occupiedSeatsAfterRelease
	if v.AvailableSeats() != availableSeats {
		t.Errorf("Expected %d available seats, got %d", availableSeats, v.AvailableSeats())
	}

	seatsToRelease = availableSeats + 1
	err = v.ReleaseSeats(seatsToRelease)
	if err == nil || err.Error() != "not enought occupied seats" {
		t.Fatalf("Expected error 'not enought occupied seats', got %v", err)
	}

	if len(v.Events()) != 1 {
		t.Errorf("Expected 1 events, got %d", len(v.Events()))
	}
	if v.Events()[0].Type() != internal.VehicleSeatsReleasedEventType {
		t.Errorf("Expected %s event, got %s", internal.VehicleSeatsReleasedEventType, v.Events()[1].Type())
	}
}

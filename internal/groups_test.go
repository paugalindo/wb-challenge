package internal_test

import (
	"testing"
	"wb-challenge/internal"
)

func TestNewGroup(t *testing.T) {
	groupID, people := 1, 3
	g, err := internal.NewGroup(groupID, people)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if g.ID() != groupID {
		t.Errorf("Expected ID to be %d, got %d", groupID, g.ID())
	}
	if g.TotalPeople() != 3 {
		t.Errorf("Expected total people to be %d, got %d", people, g.TotalPeople())
	}
	if len(g.Events()) != 1 {
		t.Errorf("Expected 1 event, got %d", len(g.Events()))
	}
	if g.Events()[0].Type() != internal.GroupCreatedEventType {
		t.Errorf("Expected %s event, got %s", internal.GroupCreatedEventType, g.Events()[0].Type())
	}

	_, err = internal.NewGroup(0, people)
	if err == nil || err.Error() != "wrong ID" {
		t.Errorf("Expected error 'wrong ID', got %v", err)
	}

	_, err = internal.NewGroup(groupID, 0)
	if err == nil || err.Error() != "wrong people number" {
		t.Errorf("Expected error 'wrong people number', got %v", err)
	}
	_, err = internal.NewGroup(groupID, 7)
	if err == nil || err.Error() != "wrong people number" {
		t.Errorf("Expected error 'wrong people number', got %v", err)
	}
}

func TestHydrateGroup(t *testing.T) {
	groupID, people, vehicleID, droppedOff := 2, 4, 2, true
	g := internal.HydrateGroup(groupID, people, vehicleID, droppedOff)
	if g.ID() != groupID {
		t.Errorf("Expected ID to be %d, got %d", groupID, g.ID())
	}
	if g.TotalPeople() != people {
		t.Errorf("Expected total people to be %d, got %d", people, g.TotalPeople())
	}
	if g.VehicleAssigned() != vehicleID {
		t.Errorf("Expected vehicle assigned to be %d, got %d", vehicleID, g.VehicleAssigned())
	}
	if g.IsDroppedOff() != droppedOff {
		t.Errorf("Expected group to be dropped off equal %t", droppedOff)
	}
}

func TestGroup_AssignVehicle(t *testing.T) {
	groupID, people := 3, 3
	g, _ := internal.NewGroup(groupID, people)

	vehicleID := 2
	err := g.AssignVehicle(vehicleID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if g.VehicleAssigned() != vehicleID {
		t.Errorf("Expected vehicle assigned to be %d, got %d", vehicleID, g.VehicleAssigned())
	}
	if len(g.Events()) != 2 {
		t.Errorf("Expected 2 events, got %d", len(g.Events()))
	}
	if g.Events()[1].Type() != internal.VehicleAssignedToGroupEventType {
		t.Errorf("Expected %s event, got %s", internal.VehicleAssignedToGroupEventType, g.Events()[1].Type())
	}

	err = g.AssignVehicle(0)
	if err == nil || err.Error() != "vehicle ID cannot be empty" {
		t.Errorf("Expected error 'vehicle ID cannot be empty', got %v", err)
	}

	err = g.AssignVehicle(vehicleID)
	if err == nil || err.Error() != "vehicle already assigned" {
		t.Errorf("Expected error 'vehicle already assigned', got %v", err)
	}

	g, _ = internal.NewGroup(groupID, people)
	g.DropOff()
	err = g.AssignVehicle(vehicleID)
	if err == nil || err.Error() != "group already dropped off" {
		t.Errorf("Expected error 'group already dropped off', got %v", err)
	}
}

func TestGroup_DropOff(t *testing.T) {
	groupID, people := 3, 3
	g, _ := internal.NewGroup(groupID, people)

	err := g.DropOff()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !g.IsDroppedOff() {
		t.Errorf("Expected group to be dropped off")
	}
	if len(g.Events()) != 2 {
		t.Errorf("Expected 2 events, got %d", len(g.Events()))
	}
	if g.Events()[1].Type() != internal.GroupDroppedOffEventType {
		t.Errorf("Expected %s event, got %s", internal.GroupDroppedOffEventType, g.Events()[1].Type())
	}

	err = g.DropOff()
	if err == nil || err.Error() != "group already dropped off" {
		t.Errorf("Expected error 'group already dropped off', got %v", err)
	}
}

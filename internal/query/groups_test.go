package query_test

import (
	"errors"
	"testing"
	"wb-challenge/internal"
	"wb-challenge/internal/query"
)

func TestFindAssignedVehicle(t *testing.T) {
	groupID, vehicleID := 1, 1
	mock := internal.NewMockGroupView()
	mock.GetFunc = func(id int) (internal.Group, error) {
		return internal.HydrateGroup(groupID, 1, vehicleID, false), nil
	}

	qs := query.NewGroupQS(mock)
	vID, err := qs.FindAssignedVehicle(groupID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if vID != vehicleID {
		t.Errorf("Expected vehicle ID %d, got %d", vehicleID, vID)
	}
	if len(mock.GetCalls) != 1 {
		t.Errorf("Expected 1 call to GroupView.Get, got %d", len(mock.GetCalls))
	}
	if mock.GetCalls[0].ID != groupID {
		t.Errorf("Expected groupID %d in GroupView.Get call, got %d", groupID, mock.GetCalls[0].ID)
	}

	mock.GetFunc = func(id int) (internal.Group, error) {
		return internal.Group{}, internal.ErrGroupNotFound
	}
	_, err = qs.FindAssignedVehicle(groupID)
	if err == nil || !errors.Is(err, query.ErrNotFound) {
		t.Fatalf("Expected error %s, got %v", query.ErrNotFound.Error(), err)
	}

	otherErr := errors.New("other error")
	mock.GetFunc = func(id int) (internal.Group, error) {
		return internal.Group{}, otherErr
	}
	_, err = qs.FindAssignedVehicle(groupID)
	if err == nil || !errors.Is(err, otherErr) {
		t.Fatalf("Expected error %s, got %v", otherErr.Error(), err)
	}
}

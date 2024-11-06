package commands_test

import (
	"context"
	"testing"
	"wb-challenge/internal"
	"wb-challenge/internal/commands"
)

func TestCreateGroupHandler_Handle(t *testing.T) {
	groupID, groupPeople := 1, 1
	mockGroupRepo := internal.NewMockGroupRepository()
	mockGroupRepo.GetFunc = func(id int) (internal.Group, error) {
		return internal.Group{}, internal.ErrGroupNotFound
	}
	mockGroupRepo.SaveFunc = func(internal.Group) error {
		return nil
	}
	mockPublisher := internal.NewMockEventsPublisher()
	mockPublisher.PublishFunc = func(...internal.Event) error {
		return nil
	}

	handler := commands.NewCreateGroupHandler(
		mockGroupRepo,
		mockPublisher,
	)

	cmd := commands.CreateGroupCmd{ID: groupID, People: groupPeople}

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

	if len(mockGroupRepo.SaveCalls) != 1 {
		t.Errorf("Expected 1 call to GroupRepository.Save, got %d", len(mockGroupRepo.SaveCalls))
	}
	if mockGroupRepo.SaveCalls[0].Group.ID() != groupID {
		t.Errorf("Expected group ID %d, got %d", groupID, mockGroupRepo.SaveCalls[0].Group.ID())
	}

	if len(mockPublisher.PublishCalls) != 1 {
		t.Errorf("Expected 1 call to EventPublisher.Publish, got %d", len(mockPublisher.PublishCalls))
	}
	if len(mockPublisher.PublishCalls[0].Events) != 1 {
		t.Errorf("Expected 1 events published, got %d", len(mockPublisher.PublishCalls[0].Events))
	}
	if mockPublisher.PublishCalls[0].Events[0].Type() != internal.GroupCreatedEventType {
		t.Errorf("Expected %s event, got %s", internal.GroupCreatedEventType, mockPublisher.PublishCalls[0].Events[0].Type())
	}
}

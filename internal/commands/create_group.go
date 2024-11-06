package commands

import (
	"context"
	"errors"
	"wb-challenge/internal"
)

const CreateGroupType = "commands.CreateGroupCmd"

type CreateGroupCmd struct {
	ID     int
	People int
}

type CreateGroupHandler struct {
	groupRepository internal.GroupRepository
	publisher       internal.EventsPublisher
}

func NewCreateGroupHandler(groupRepository internal.GroupRepository, publisher internal.EventsPublisher) CreateGroupHandler {
	return CreateGroupHandler{
		groupRepository: groupRepository,
		publisher:       publisher,
	}
}

var ErrGroupAlreadyExist = errors.New("group already exists")

func (h *CreateGroupHandler) Handle(ctx context.Context, c any) error {
	cmd, ok := c.(CreateGroupCmd)
	if !ok {
		return errors.New("invalid command type")
	}

	_, err := h.groupRepository.Get(cmd.ID)
	if err != nil && !errors.Is(err, internal.ErrGroupNotFound) {
		return err
	}
	if err == nil {
		return ErrGroupAlreadyExist
	}

	g, err := internal.NewGroup(cmd.ID, cmd.People)
	if err != nil {
		return err
	}

	if err = h.groupRepository.Save(g); err != nil {
		return err
	}

	return h.publisher.Publish(g.Events()...)
}

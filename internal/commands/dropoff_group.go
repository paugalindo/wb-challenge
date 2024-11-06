package commands

import (
	"context"
	"errors"
	"sync"
	"wb-challenge/internal"
)

const DropOffGroupType = "commands.DropOffGroupCmd"

type DropOffGroupCmd struct {
	ID int
}

type DropOffGroupHandler struct {
	mutex           *sync.Mutex
	groupRepository internal.GroupRepository
	publisher       internal.EventsPublisher
}

func NewDropOffGroupHandler(mutex *sync.Mutex, groupRepository internal.GroupRepository,
	publisher internal.EventsPublisher,
) DropOffGroupHandler {
	return DropOffGroupHandler{
		mutex:           mutex,
		groupRepository: groupRepository,
		publisher:       publisher,
	}
}

var ErrNotFound = errors.New("not found")

func (h *DropOffGroupHandler) Handle(ctx context.Context, c any) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	cmd, ok := c.(DropOffGroupCmd)
	if !ok {
		return errors.New("invalid command type")
	}

	g, err := h.groupRepository.Get(cmd.ID)
	if err != nil {
		if errors.Is(err, internal.ErrGroupNotFound) {
			return ErrNotFound
		}
		return err
	}

	if err = g.DropOff(); err != nil {
		return err
	}

	if err = h.groupRepository.Save(g); err != nil {
		return err
	}

	return h.publisher.Publish(g.Events()...)
}

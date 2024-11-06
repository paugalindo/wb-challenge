package bus

import (
	"context"
	"errors"
	"fmt"
)

type CommandHandler interface {
	Handle(context.Context, any) error
}

type CommandBus struct {
	handlers map[string]CommandHandler
}

func NewCommandBus() CommandBus {
	return CommandBus{
		handlers: make(map[string]CommandHandler),
	}
}

func (c *CommandBus) RegisterHandler(cmdType string, h CommandHandler) {
	c.handlers[cmdType] = h
}

func (c CommandBus) Dispatch(ctx context.Context, cmd any) error {
	cmdType := fmt.Sprintf("%T", cmd)
	handler, ok := c.handlers[cmdType]
	if !ok {
		return errors.New("no handler registered for command type")
	}
	return handler.Handle(ctx, cmd)
}

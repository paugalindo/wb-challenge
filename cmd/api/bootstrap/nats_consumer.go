package bootstrap

import (
	"log"
	"wb-challenge/bus"
	"wb-challenge/internal"
	eventhandlers "wb-challenge/internal/event-handlers"
	"wb-challenge/internal/nats"

	natsconn "github.com/nats-io/nats.go"
)

func InitNATSConsumer(conn *natsconn.Conn, cmdBus *bus.CommandBus, logger *log.Logger) error {
	consumer := nats.NewConsumer(conn)

	groupCreatedHandler := eventhandlers.NewGroupCreated(cmdBus, logger)
	consumer.Subscribe(internal.GroupCreatedEventType, groupCreatedHandler.Handle)

	vehicleAssignedToGroupHandler := eventhandlers.NewVehicleAssignedToGroup(cmdBus, logger)
	consumer.Subscribe(internal.VehicleAssignedToGroupEventType, vehicleAssignedToGroupHandler.Handle)

	groupDroppedOffHandler := eventhandlers.NewGroupDroppedOff(cmdBus, logger)
	consumer.Subscribe(internal.GroupDroppedOffEventType, groupDroppedOffHandler.Handle)

	vehicleSeatsReleasedHandler := eventhandlers.NewVehicleSeatsReleased(cmdBus, logger)
	consumer.Subscribe(internal.VehicleSeatsReleasedEventType, vehicleSeatsReleasedHandler.Handle)

	return nil
}

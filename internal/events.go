package internal

type EventsPublisher interface {
	Publish(...Event) error
}

type Event interface {
	Type() string
}

// Vehicle events

type VehicleCreatedEvent struct {
	ID    int
	Seats int
}

const VehicleCreatedEventType = "vehicle.created"

func (e VehicleCreatedEvent) Type() string {
	return VehicleCreatedEventType
}

type VehicleSeatsOccupiedEvent struct {
	ID    int
	Seats int
}

const VehicleSeatsOccupiedEventType = "vehicle.occupied_seats"

func (e VehicleSeatsOccupiedEvent) Type() string {
	return VehicleSeatsOccupiedEventType
}

type VehicleSeatsReleasedEvent struct {
	ID    int
	Seats int
}

const VehicleSeatsReleasedEventType = "vehicle.released_seats"

func (e VehicleSeatsReleasedEvent) Type() string {
	return VehicleSeatsReleasedEventType
}

// Group events

type GroupCreatedEvent struct {
	ID     int
	People int
}

const GroupCreatedEventType = "group.created"

func (e GroupCreatedEvent) Type() string {
	return GroupCreatedEventType
}

type VehicleAssignedToGroupEvent struct {
	ID        int
	VehicleID int
}

const VehicleAssignedToGroupEventType = "group.vehicle_assigned"

func (e VehicleAssignedToGroupEvent) Type() string {
	return VehicleAssignedToGroupEventType
}

type GroupDroppedOffEvent struct {
	ID int
}

const GroupDroppedOffEventType = "group.dropped_off"

func (e GroupDroppedOffEvent) Type() string {
	return GroupDroppedOffEventType
}

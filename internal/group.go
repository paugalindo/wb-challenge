package internal

import "errors"

type GroupRepository interface {
	Save(Group) error
	RemoveAllGroups() error
	Get(id int) (Group, error)
	GetUnassignedOrderedByCreatedAt() ([]Group, error)
}

type GroupView interface {
	Get(id int) (Group, error)
}

var ErrGroupNotFound = errors.New("group not found")

type Group struct {
	id              int
	people          []People
	vehicleAssigned int
	droppedOff      bool

	events []Event
}

func NewGroup(id, people int) (Group, error) {
	g := Group{
		id:     id,
		people: make([]People, people),
	}
	if err := g.validate(); err != nil {
		return Group{}, err
	}

	g.events = append(g.events, GroupCreatedEvent{
		ID:     id,
		People: people,
	})
	return g, nil
}

const (
	minPeopleNumber = 1
	maxPeopleNumber = 6
)

func (g *Group) validate() error {
	if g.id == 0 {
		return errors.New("wrong ID")
	}
	if len(g.people) < minPeopleNumber || len(g.people) > maxPeopleNumber {
		return errors.New("wrong people number")
	}
	return nil
}

func HydrateGroup(id, people, vehicleID int, droppedOff bool) Group {
	return Group{
		id:              id,
		people:          make([]People, people),
		vehicleAssigned: vehicleID,
		droppedOff:      droppedOff,
	}
}

func (g *Group) ID() int {
	return g.id
}

func (g *Group) TotalPeople() int {
	return len(g.people)
}

func (g *Group) VehicleAssigned() int {
	return g.vehicleAssigned
}

func (g *Group) IsDroppedOff() bool {
	return g.droppedOff
}

func (g *Group) AssignVehicle(vehicleID int) error {
	if vehicleID == 0 {
		return errors.New("vehicle ID cannot be empty")
	}
	if g.vehicleAssigned != 0 {
		return errors.New("vehicle already assigned")
	}
	if g.droppedOff {
		return errors.New("group already dropped off")
	}

	g.vehicleAssigned = vehicleID

	g.events = append(g.events, VehicleAssignedToGroupEvent{
		ID:        g.ID(),
		VehicleID: vehicleID,
	})
	return nil
}

func (g *Group) DropOff() error {
	if g.droppedOff {
		return errors.New("group already dropped off")
	}

	g.droppedOff = true

	g.events = append(g.events, GroupDroppedOffEvent{
		ID: g.ID(),
	})
	return nil
}

func (g *Group) Events() []Event {
	return g.events
}

type People struct {
}

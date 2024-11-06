package internal

import "errors"

type VehicleRepository interface {
	Save(Vehicle) error
	RemoveAllVehicles() error
	Get(id int) (Vehicle, error)
}

type VehicleView interface {
	GetWithEmptySeats(seats int) (Vehicle, error)
}

var ErrVehicleNotFound = errors.New("vehicle not found")

type Vehicle struct {
	id    int
	seats []Seat

	events []Event
}

func NewVehicle(id, seats int) (Vehicle, error) {
	v := Vehicle{
		id:    id,
		seats: make([]Seat, seats),
	}
	if err := v.validate(); err != nil {
		return Vehicle{}, err
	}

	v.events = append(v.events, VehicleCreatedEvent{
		ID:    id,
		Seats: seats,
	})
	return v, nil
}

const (
	minSeatsNumber = 4
	maxSeatsNumber = 6
)

func (v *Vehicle) validate() error {
	if v.id == 0 {
		return errors.New("wrong ID")
	}
	if len(v.seats) < minSeatsNumber || len(v.seats) > maxSeatsNumber {
		return errors.New("wrong number of seats")
	}
	return nil
}

func HydrateVehicle(id, totalSeats, occupiedSeats int) Vehicle {
	v := Vehicle{
		id:    id,
		seats: make([]Seat, totalSeats),
	}

	for i := 0; i < occupiedSeats; i++ {
		v.seats[i].occupied = true
	}

	return v
}

func (v *Vehicle) ID() int {
	return v.id
}

func (v *Vehicle) AvailableSeats() int {
	available := 0
	for _, s := range v.seats {
		if !s.occupied {
			available++
		}
	}
	return available
}

func (v *Vehicle) OccupiedSeats() int {
	occupied := 0
	for _, s := range v.seats {
		if s.occupied {
			occupied++
		}
	}
	return occupied
}

func (v *Vehicle) OccupySeats(s int) error {
	if v.AvailableSeats() < s {
		return errors.New("not enought available seats")
	}

	for i := 0; i < s; i++ {
		for i, s := range v.seats {
			if !s.occupied {
				v.seats[i] = s.Occupy()
				break
			}
		}
	}

	v.events = append(v.events, VehicleSeatsOccupiedEvent{
		ID:    v.ID(),
		Seats: s,
	})
	return nil
}

func (v *Vehicle) ReleaseSeats(s int) error {
	if v.OccupiedSeats() < s {
		return errors.New("not enought occupied seats")
	}

	for i := 0; i < s; i++ {
		for i, s := range v.seats {
			if s.occupied {
				v.seats[i] = s.Release()
				break
			}
		}
	}

	v.events = append(v.events, VehicleSeatsReleasedEvent{
		ID:    v.ID(),
		Seats: s,
	})
	return nil
}

func (v *Vehicle) Events() []Event {
	return v.events
}

type Seat struct {
	occupied bool
}

func NewSeat() Seat {
	return Seat{}
}

func (s *Seat) Occupy() Seat {
	seat := NewSeat()
	seat.occupied = true
	return seat
}

func (s *Seat) Release() Seat {
	return NewSeat()
}

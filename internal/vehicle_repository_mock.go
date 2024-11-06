package internal

type MockVehicleRepository struct {
	SaveFunc              func(Vehicle) error
	RemoveAllFunc         func() error
	GetFunc               func(id int) (Vehicle, error)
	GetWithEmptySeatsFunc func(seats int) (Vehicle, error)

	SaveCalls              []SaveCall
	RemoveAllCalls         []RemoveAllCall
	GetCalls               []GetCall
	GetWithEmptySeatsCalls []GetWithEmptySeatsCall
}

type SaveCall struct {
	Vehicle Vehicle
}

type RemoveAllCall struct{}

type GetCall struct {
	ID int
}

type GetWithEmptySeatsCall struct {
	Seats int
}

func NewMockVehicleRepository() *MockVehicleRepository {
	return &MockVehicleRepository{
		SaveCalls:              []SaveCall{},
		RemoveAllCalls:         []RemoveAllCall{},
		GetCalls:               []GetCall{},
		GetWithEmptySeatsCalls: []GetWithEmptySeatsCall{},
	}
}

func (m *MockVehicleRepository) Save(vehicle Vehicle) error {
	m.SaveCalls = append(m.SaveCalls, SaveCall{Vehicle: vehicle})
	if m.SaveFunc != nil {
		return m.SaveFunc(vehicle)
	}
	return nil
}

func (m *MockVehicleRepository) RemoveAllVehicles() error {
	m.RemoveAllCalls = append(m.RemoveAllCalls, RemoveAllCall{})
	if m.RemoveAllFunc != nil {
		return m.RemoveAllFunc()
	}
	return nil
}

func (m *MockVehicleRepository) Get(id int) (Vehicle, error) {
	m.GetCalls = append(m.GetCalls, GetCall{ID: id})
	if m.GetFunc != nil {
		return m.GetFunc(id)
	}
	return Vehicle{}, nil
}

func (m *MockVehicleRepository) GetWithEmptySeats(seats int) (Vehicle, error) {
	m.GetWithEmptySeatsCalls = append(m.GetWithEmptySeatsCalls, GetWithEmptySeatsCall{Seats: seats})
	if m.GetWithEmptySeatsFunc != nil {
		return m.GetWithEmptySeatsFunc(seats)
	}
	return Vehicle{}, nil
}

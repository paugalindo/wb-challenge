package internal

import (
	"errors"
)

type MockGroupRepository struct {
	SaveFunc                            func(Group) error
	RemoveAllFunc                       func() error
	GetFunc                             func(id int) (Group, error)
	GetUnassignedOrderedByCreatedAtFunc func() ([]Group, error)

	SaveCalls                            []SaveGroupCall
	RemoveAllCalls                       []RemoveAllGroupsCall
	GetCalls                             []GetGroupCall
	GetUnassignedOrderedByCreatedAtCalls []GetUnassignedOrderedByCreatedAtGroupsCall
}

type SaveGroupCall struct {
	Group Group
}

type RemoveAllGroupsCall struct{}

type GetGroupCall struct {
	ID int
}

type GetUnassignedOrderedByCreatedAtGroupsCall struct{}

func NewMockGroupRepository() *MockGroupRepository {
	return &MockGroupRepository{
		SaveCalls:                            []SaveGroupCall{},
		RemoveAllCalls:                       []RemoveAllGroupsCall{},
		GetCalls:                             []GetGroupCall{},
		GetUnassignedOrderedByCreatedAtCalls: []GetUnassignedOrderedByCreatedAtGroupsCall{},
	}
}

func (m *MockGroupRepository) Save(group Group) error {
	m.SaveCalls = append(m.SaveCalls, SaveGroupCall{Group: group})
	if m.SaveFunc != nil {
		return m.SaveFunc(group)
	}
	return nil
}

func (m *MockGroupRepository) RemoveAllGroups() error {
	m.RemoveAllCalls = append(m.RemoveAllCalls, RemoveAllGroupsCall{})
	if m.RemoveAllFunc != nil {
		return m.RemoveAllFunc()
	}
	return nil
}

func (m *MockGroupRepository) Get(id int) (Group, error) {
	m.GetCalls = append(m.GetCalls, GetGroupCall{ID: id})
	if m.GetFunc != nil {
		return m.GetFunc(id)
	}
	return Group{}, errors.New("group not found")
}

func (m *MockGroupRepository) GetUnassignedOrderedByCreatedAt() ([]Group, error) {
	m.GetUnassignedOrderedByCreatedAtCalls = append(m.GetUnassignedOrderedByCreatedAtCalls, GetUnassignedOrderedByCreatedAtGroupsCall{})
	if m.GetUnassignedOrderedByCreatedAtFunc != nil {
		return m.GetUnassignedOrderedByCreatedAtFunc()
	}
	return nil, errors.New("no unassigned groups found")
}

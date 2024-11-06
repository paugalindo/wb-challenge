package internal

import (
	"errors"
)

type MockGroupView struct {
	GetFunc func(id int) (Group, error)

	GetCalls []GetGroupCall
}

func NewMockGroupView() *MockGroupView {
	return &MockGroupView{
		GetCalls: []GetGroupCall{},
	}
}

func (m *MockGroupView) Get(id int) (Group, error) {
	m.GetCalls = append(m.GetCalls, GetGroupCall{ID: id})
	if m.GetFunc != nil {
		return m.GetFunc(id)
	}
	return Group{}, errors.New("group not found")
}

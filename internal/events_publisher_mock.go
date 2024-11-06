package internal

type MockEventsPublisher struct {
	PublishFunc func(...Event) error

	PublishCalls []PublishCall
}

type PublishCall struct {
	Events []Event
}

func NewMockEventsPublisher() *MockEventsPublisher {
	return &MockEventsPublisher{
		PublishCalls: []PublishCall{},
	}
}

func (m *MockEventsPublisher) Publish(events ...Event) error {
	m.PublishCalls = append(m.PublishCalls, PublishCall{Events: events})
	if m.PublishCalls != nil {
		return m.PublishFunc(events...)
	}
	return nil
}

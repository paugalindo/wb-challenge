package query

import (
	"errors"
	"wb-challenge/internal"
)

type GroupQS struct {
	groupView internal.GroupView
}

func NewGroupQS(groupView internal.GroupView) GroupQS {
	return GroupQS{
		groupView: groupView,
	}
}

var ErrNotFound = errors.New("not found")

func (sq GroupQS) FindAssignedVehicle(id int) (int, error) {
	g, err := sq.groupView.Get(id)
	if err != nil {
		if errors.Is(err, internal.ErrGroupNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	return g.VehicleAssigned(), nil
}

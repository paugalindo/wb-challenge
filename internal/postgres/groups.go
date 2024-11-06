package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"wb-challenge/internal"
)

type GroupsRepository struct {
	sql *sql.DB
}

func NewGroupsRepository(sql *sql.DB) GroupsRepository {
	return GroupsRepository{sql: sql}
}

const insertGroupSmt = `INSERT INTO groups (id, vehicle_id, value) VALUES ($1, $2, $3) ON CONFLICT(id) DO UPDATE SET vehicle_id=EXCLUDED.vehicle_id, value=EXCLUDED.value`

func (r *GroupsRepository) Save(g internal.Group) error {
	dto := groupToDTO(g)
	value, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	_, err = r.sql.Exec(insertGroupSmt, dto.ID, dto.VehicleAssigned, value)
	return err
}

const removeAllGroupsSmt = `TRUNCATE TABLE groups CASCADE`

func (r *GroupsRepository) RemoveAllGroups() error {
	_, err := r.sql.Exec(removeAllGroupsSmt)
	return err
}

const selectGroupByIDSmt = `SELECT value FROM groups WHERE id = $1`

func (r *GroupsRepository) Get(id int) (internal.Group, error) {
	var rawValue string
	if err := r.sql.QueryRow(selectGroupByIDSmt, id).Scan(&rawValue); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.Group{}, internal.ErrGroupNotFound
		}
		return internal.Group{}, err
	}

	dto := GroupDTO{}
	if err := json.Unmarshal([]byte(rawValue), &dto); err != nil {
		return internal.Group{}, err
	}

	return groupFromDTO(dto), nil
}

const selectUnassignedGroupsSmt = `
SELECT value FROM groups 
WHERE (value->>'vehicle_assigned')::int = 0 AND (value->>'dropped_off')::boolean = FALSE
ORDER BY created_at;`

func (r *GroupsRepository) GetUnassignedOrderedByCreatedAt() ([]internal.Group, error) {
	rows, err := r.sql.Query(selectUnassignedGroupsSmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]internal.Group, 0)
	for rows.Next() {
		var rawValue string
		if err := rows.Scan(&rawValue); err != nil {
			return nil, err
		}
		dto := GroupDTO{}
		if err := json.Unmarshal([]byte(rawValue), &dto); err != nil {
			return nil, err
		}
		groups = append(groups, groupFromDTO(dto))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

type GroupDTO struct {
	ID              int         `json:"id"`
	People          []PeopleDTO `json:"people"`
	VehicleAssigned int         `json:"vehicle_assigned"`
	DroppedOff      bool        `json:"dropped_off"`
}

type PeopleDTO struct{}

func groupToDTO(group internal.Group) GroupDTO {
	return GroupDTO{
		ID:              group.ID(),
		People:          make([]PeopleDTO, group.TotalPeople()),
		VehicleAssigned: group.VehicleAssigned(),
		DroppedOff:      group.IsDroppedOff(),
	}
}

func groupFromDTO(dto GroupDTO) internal.Group {
	return internal.HydrateGroup(dto.ID, len(dto.People), dto.VehicleAssigned, dto.DroppedOff)
}

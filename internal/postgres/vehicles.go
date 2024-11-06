package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"wb-challenge/internal"
)

type VehiclesRepository struct {
	sql *sql.DB
}

func NewVehiclesRepository(sql *sql.DB) VehiclesRepository {
	return VehiclesRepository{sql: sql}
}

const insertVehicleSmt = `INSERT INTO vehicles (id, value) VALUES ($1, $2) ON CONFLICT(id) DO UPDATE SET value=EXCLUDED.value`

func (r *VehiclesRepository) Save(v internal.Vehicle) error {
	dto := vehicleToDTO(v)
	value, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	_, err = r.sql.Exec(insertVehicleSmt, dto.ID, value)
	return err
}

const removeAllVehiclesSmt = `TRUNCATE TABLE vehicles CASCADE`

func (r *VehiclesRepository) RemoveAllVehicles() error {
	_, err := r.sql.Exec(removeAllVehiclesSmt)
	return err
}

const selectVechileByIDSmt = `SELECT value FROM vehicles WHERE id = $1`

func (r *VehiclesRepository) Get(id int) (internal.Vehicle, error) {
	var rawValue string
	if err := r.sql.QueryRow(selectVechileByIDSmt, id).Scan(&rawValue); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.Vehicle{}, internal.ErrVehicleNotFound
		}
		return internal.Vehicle{}, err
	}

	dto := VehicleDTO{}
	if err := json.Unmarshal([]byte(rawValue), &dto); err != nil {
		return internal.Vehicle{}, err
	}

	return vehicleFromDTO(dto), nil
}

const selectVehicleWithEmptySeatsSmt = `
SELECT value FROM vehicles 
WHERE (
    SELECT COUNT(*) 
    FROM jsonb_array_elements(value->'seats') AS seat
    WHERE (seat->>'occupied')::boolean = FALSE
) >= $1 AND (
    SELECT COUNT(*)
    FROM jsonb_array_elements(value->'seats') AS seat
    WHERE (seat->>'occupied')::boolean = TRUE
) = 0
ORDER BY (
    SELECT COUNT(*)
    FROM jsonb_array_elements(value->'seats') AS seat
    WHERE (seat->>'occupied')::boolean = FALSE
) ASC;`

func (r *VehiclesRepository) GetWithEmptySeats(seats int) (internal.Vehicle, error) {
	var rawValue string
	if err := r.sql.QueryRow(selectVehicleWithEmptySeatsSmt, seats).Scan(&rawValue); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.Vehicle{}, internal.ErrVehicleNotFound
		}
		return internal.Vehicle{}, err
	}

	dto := VehicleDTO{}
	if err := json.Unmarshal([]byte(rawValue), &dto); err != nil {
		return internal.Vehicle{}, err
	}

	return vehicleFromDTO(dto), nil
}

type VehicleDTO struct {
	ID    int       `json:"id"`
	Seats []SeatDTO `json:"seats"`
}

type SeatDTO struct {
	Occupied bool `json:"occupied"`
}

func vehicleToDTO(vehicle internal.Vehicle) VehicleDTO {
	dto := VehicleDTO{
		ID:    vehicle.ID(),
		Seats: make([]SeatDTO, vehicle.AvailableSeats()+vehicle.OccupiedSeats()),
	}

	for i := 0; i < vehicle.OccupiedSeats(); i++ {
		dto.Seats[i].Occupied = true
	}

	return dto
}

func vehicleFromDTO(dto VehicleDTO) internal.Vehicle {
	occupiedSeats := 0
	for _, s := range dto.Seats {
		if s.Occupied {
			occupiedSeats++
		}
	}

	return internal.HydrateVehicle(dto.ID, len(dto.Seats), occupiedSeats)
}

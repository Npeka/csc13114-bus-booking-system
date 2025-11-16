package repository

import (
	"gorm.io/gorm"
)

type Repositories struct {
	Trip     TripRepository
	Route    RouteRepository
	Bus      BusRepository
	Operator OperatorRepository
	Seat     SeatRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Trip:     NewTripRepository(db),
		Route:    NewRouteRepository(db),
		Bus:      NewBusRepository(db),
		Operator: NewOperatorRepository(db),
		Seat:     NewSeatRepository(db),
	}
}

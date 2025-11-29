package service

import (
	"context"

	"bus-booking/trip-service/internal/constants"
	"bus-booking/trip-service/internal/model"
)

type ConstantsService interface {
	GetBusConstants(ctx context.Context) (*model.BusConstants, error)
	GetRouteConstants(ctx context.Context) (*model.RouteConstants, error)
	GetTripConstants(ctx context.Context) (*model.TripConstants, error)
	GetAllConstants(ctx context.Context) (*model.ConstantsResponse, error)
}

type ConstantsServiceImpl struct{}

func NewConstantsService() ConstantsService {
	return &ConstantsServiceImpl{}
}

func (s *ConstantsServiceImpl) GetBusConstants(ctx context.Context) (*model.BusConstants, error) {
	seatTypes := make([]model.SeatTypeConstant, 0)
	for _, st := range constants.AllSeatTypes() {
		seatTypes = append(seatTypes, model.SeatTypeConstant{
			Value:           st.String(),
			DisplayName:     getSeatTypeDisplayName(st),
			PriceMultiplier: st.GetPriceMultiplier(),
		})
	}

	amenities := make([]model.AmenityConstant, 0)
	for _, a := range constants.AllAmenities {
		amenities = append(amenities, model.AmenityConstant{
			Value:       a.String(),
			DisplayName: a.GetDisplayName(),
		})
	}

	busTypes := make([]model.BusTypeConstant, 0)
	for _, bt := range constants.AllBusTypes() {
		busTypes = append(busTypes, model.BusTypeConstant{
			Value:       bt.String(),
			DisplayName: getBusTypeDisplayName(bt),
		})
	}

	return &model.BusConstants{
		SeatTypes: seatTypes,
		Amenities: amenities,
		BusTypes:  busTypes,
	}, nil
}

func (s *ConstantsServiceImpl) GetRouteConstants(ctx context.Context) (*model.RouteConstants, error) {
	stopTypes := make([]model.StopTypeConstant, 0)
	for _, st := range constants.AllStopTypes() {
		stopTypes = append(stopTypes, model.StopTypeConstant{
			Value:       st.String(),
			DisplayName: getStopTypeDisplayName(st),
		})
	}

	return &model.RouteConstants{
		StopTypes: stopTypes,
	}, nil
}

func (s *ConstantsServiceImpl) GetTripConstants(ctx context.Context) (*model.TripConstants, error) {
	tripStatuses := make([]model.TripStatusConstant, 0)
	for _, ts := range constants.AllTripStatuses() {
		tripStatuses = append(tripStatuses, model.TripStatusConstant{
			Value:       ts.String(),
			DisplayName: getTripStatusDisplayName(ts),
		})
	}

	return &model.TripConstants{
		TripStatuses: tripStatuses,
	}, nil
}

func (s *ConstantsServiceImpl) GetAllConstants(ctx context.Context) (*model.ConstantsResponse, error) {
	busConstants, err := s.GetBusConstants(ctx)
	if err != nil {
		return nil, err
	}

	routeConstants, err := s.GetRouteConstants(ctx)
	if err != nil {
		return nil, err
	}

	tripConstants, err := s.GetTripConstants(ctx)
	if err != nil {
		return nil, err
	}

	return &model.ConstantsResponse{
		Bus:   *busConstants,
		Route: *routeConstants,
		Trip:  *tripConstants,
	}, nil
}

func getSeatTypeDisplayName(st constants.SeatType) string {
	switch st {
	case constants.SeatTypeStandard:
		return "Standard"
	case constants.SeatTypeVIP:
		return "VIP"
	case constants.SeatTypeSleeper:
		return "Sleeper"
	default:
		return st.String()
	}
}

func getBusTypeDisplayName(bt constants.BusType) string {
	switch bt {
	case constants.BusTypeStandard:
		return "Standard"
	case constants.BusTypeVIP:
		return "VIP"
	case constants.BusTypeSleeper:
		return "Sleeper"
	case constants.BusTypeDoubleDecker:
		return "Double Decker"
	default:
		return bt.String()
	}
}

func getStopTypeDisplayName(st constants.StopType) string {
	switch st {
	case constants.StopTypePickup:
		return "Pickup Only"
	case constants.StopTypeDropoff:
		return "Dropoff Only"
	case constants.StopTypeBoth:
		return "Pickup & Dropoff"
	default:
		return st.String()
	}
}

func getTripStatusDisplayName(ts constants.TripStatus) string {
	switch ts {
	case constants.TripStatusScheduled:
		return "Scheduled"
	case constants.TripStatusInProgress:
		return "In Progress"
	case constants.TripStatusCompleted:
		return "Completed"
	case constants.TripStatusCancelled:
		return "Cancelled"
	case constants.TripStatusDelayed:
		return "Delayed"
	default:
		return ts.String()
	}
}

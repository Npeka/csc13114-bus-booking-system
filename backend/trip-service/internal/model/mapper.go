package model

import (
	"bus-booking/trip-service/internal/constants"
)

// ToBusResponse converts Bus entity to BusResponse with mapped constants
func ToBusResponse(bus *Bus) *BusResponse {
	if bus == nil {
		return nil
	}

	// Map amenities to ConstantDisplay
	amenities := make([]ConstantDisplay, len(bus.Amenities))
	for i, a := range bus.Amenities {
		amenity := constants.Amenity(a)
		amenities[i] = ConstantDisplay{
			Value:       amenity.String(),
			DisplayName: amenity.GetDisplayName(),
		}
	}

	// Map seats if present
	var seats []SeatResponse
	if len(bus.Seats) > 0 {
		seats = make([]SeatResponse, len(bus.Seats))
		for i, seat := range bus.Seats {
			seats[i] = *ToSeatResponse(&seat)
		}
	}

	return &BusResponse{
		ID:           bus.ID,
		PlateNumber:  bus.PlateNumber,
		Model:        bus.Model,
		SeatCapacity: bus.SeatCapacity,
		Amenities:    amenities,
		IsActive:     bus.IsActive,
		Seats:        seats,
	}
}

// ToSeatResponse converts Seat entity to SeatResponse with mapped constants
func ToSeatResponse(seat *Seat) *SeatResponse {
	if seat == nil {
		return nil
	}

	return &SeatResponse{
		ID:         seat.ID,
		BusID:      seat.BusID,
		SeatNumber: seat.SeatNumber,
		Row:        seat.Row,
		Column:     seat.Column,
		SeatType: ConstantDisplay{
			Value:       seat.SeatType.String(),
			DisplayName: seat.SeatType.GetDisplayName(),
		},
		PriceMultiplier: seat.PriceMultiplier,
		IsAvailable:     seat.IsAvailable,
		Floor:           seat.Floor,
	}
}

// ToRouteStopResponse converts RouteStop entity to RouteStopResponse with mapped constants
func ToRouteStopResponse(stop *RouteStop) *RouteStopResponse {
	if stop == nil {
		return nil
	}

	return &RouteStopResponse{
		ID:        stop.ID,
		RouteID:   stop.RouteID,
		StopOrder: stop.StopOrder,
		StopType: ConstantDisplay{
			Value:       stop.StopType.String(),
			DisplayName: stop.StopType.GetDisplayName(),
		},
		Location:      stop.Location,
		Address:       stop.Address,
		Latitude:      stop.Latitude,
		Longitude:     stop.Longitude,
		OffsetMinutes: stop.OffsetMinutes,
		IsActive:      stop.IsActive,
	}
}

// ToRouteResponse converts Route entity to RouteResponse with mapped constants
func ToRouteResponse(route *Route) *RouteResponse {
	if route == nil {
		return nil
	}

	// Map route stops if present
	var routeStops []RouteStopResponse
	if len(route.RouteStops) > 0 {
		routeStops = make([]RouteStopResponse, len(route.RouteStops))
		for i, stop := range route.RouteStops {
			routeStops[i] = *ToRouteStopResponse(&stop)
		}
	}

	return &RouteResponse{
		ID:               route.ID,
		CreatedAt:        route.CreatedAt,
		UpdatedAt:        route.UpdatedAt,
		Origin:           route.Origin,
		Destination:      route.Destination,
		DistanceKm:       route.DistanceKm,
		EstimatedMinutes: route.EstimatedMinutes,
		IsActive:         route.IsActive,
		RouteStops:       routeStops,
	}
}

// ToTripResponse converts Trip entity to TripResponse with mapped constants
func ToTripResponse(trip *Trip) *TripResponse {
	if trip == nil {
		return nil
	}

	status := constants.TripStatus(trip.Status)

	return &TripResponse{
		ID:            trip.ID,
		RouteID:       trip.RouteID,
		BusID:         trip.BusID,
		DepartureTime: trip.DepartureTime,
		ArrivalTime:   trip.ArrivalTime,
		BasePrice:     trip.BasePrice,
		Status: ConstantDisplay{
			Value:       status.String(),
			DisplayName: status.GetDisplayName(),
		},
		IsActive:  trip.IsActive,
		Route:     ToRouteResponse(trip.Route),
		Bus:       ToBusResponse(trip.Bus),
		CreatedAt: trip.CreatedAt,
		UpdatedAt: trip.UpdatedAt,
	}
}

// ToBusResponseList converts list of Bus entities to BusResponse list
func ToBusResponseList(buses []Bus) []BusResponse {
	responses := make([]BusResponse, len(buses))
	for i, bus := range buses {
		responses[i] = *ToBusResponse(&bus)
	}
	return responses
}

// ToSeatResponseList converts list of Seat entities to SeatResponse list
func ToSeatResponseList(seats []Seat) []SeatResponse {
	responses := make([]SeatResponse, len(seats))
	for i, seat := range seats {
		responses[i] = *ToSeatResponse(&seat)
	}
	return responses
}

// ToRouteResponseList converts list of Route entities to RouteResponse list
func ToRouteResponseList(routes []Route) []RouteResponse {
	responses := make([]RouteResponse, len(routes))
	for i, route := range routes {
		responses[i] = *ToRouteResponse(&route)
	}
	return responses
}

// ToTripResponseList converts list of Trip entities to TripResponse list
func ToTripResponseList(trips []Trip) []TripResponse {
	responses := make([]TripResponse, len(trips))
	for i, trip := range trips {
		responses[i] = *ToTripResponse(&trip)
	}
	return responses
}

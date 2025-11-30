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
	GetSearchFilterConstants(ctx context.Context) (*model.SearchFilterConstants, error)
	GetCities(ctx context.Context) ([]string, error)
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
			DisplayName:     st.GetDisplayName(),
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
			DisplayName: bt.GetDisplayName(),
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
			DisplayName: st.GetDisplayName(),
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
			DisplayName: ts.GetDisplayName(),
		})
	}

	return &model.TripConstants{
		TripStatuses: tripStatuses,
	}, nil
}

func (s *ConstantsServiceImpl) GetSearchFilterConstants(ctx context.Context) (*model.SearchFilterConstants, error) {
	// Price ranges - single range for user selection
	priceRanges := []model.FilterPriceRange{
		{Min: 0, Max: 1000000},
	}

	// Time slots
	timeSlots := []model.FilterTimeSlot{
		{StartTime: "00:00", EndTime: "06:00", DisplayName: "Sáng sớtm (00:00 - 06:00)"},
		{StartTime: "06:00", EndTime: "12:00", DisplayName: "Ban ngày (06:00 - 12:00)"},
		{StartTime: "12:00", EndTime: "18:00", DisplayName: "Chiều (12:00 - 18:00)"},
		{StartTime: "18:00", EndTime: "24:00", DisplayName: "Tối (18:00 - 24:00)"},
	}

	// Seat types - reuse from bus constants
	seatTypes := make([]model.SeatTypeConstant, 0)
	for _, st := range constants.AllSeatTypes() {
		seatTypes = append(seatTypes, model.SeatTypeConstant{
			Value:           st.String(),
			DisplayName:     st.GetDisplayName(),
			PriceMultiplier: st.GetPriceMultiplier(),
		})
	}

	// Amenities - filter common ones for search
	commonAmenities := []constants.Amenity{
		constants.AmenityWiFi,
		constants.AmenityAC,
		constants.AmenityWater,
		constants.AmenityCharging,
		constants.AmenityToilet,
		constants.AmenityTV,
	}
	amenities := make([]model.AmenityConstant, 0)
	for _, a := range commonAmenities {
		amenities = append(amenities, model.AmenityConstant{
			Value:       a.String(),
			DisplayName: a.GetDisplayName(),
		})
	}

	// Sort options
	sortOptions := []model.ConstantDisplay{
		{Value: "departure_time", DisplayName: "Giờ đi"},
		{Value: "price", DisplayName: "Giá vé"},
		{Value: "duration", DisplayName: "Thời gian hành trình"},
	}

	// Vietnam cities
	cities := []string{
		"Hà Nội",
		"TP. Hồ Chí Minh",
		"Đà Nẵng",
		"Đà Lạt",
		"Hải Phòng",
		"Cần Thơ",
		"Huế",
		"Nha Trang",
		"Quảng Ninh",
		"Bắc Ninh",
		"Bình Dương",
		"Bình Định",
		"Bình Phước",
		"Bến Tre",
		"Gia Lai",
		"Kiên Giang",
		"Lâm Đồng",
		"Long An",
		"Nam Định",
		"Nghệ An",
		"Phú Yên",
		"Quảng Bình",
		"Quảng Nam",
		"Quảng Ngãi",
		"Sóc Trăng",
		"Tây Ninh",
		"Thái Nguyên",
		"Thanh Hóa",
		"Tiền Giang",
		"Trà Vinh",
		"Vĩnh Long",
		"Vũng Tàu",
	}

	return &model.SearchFilterConstants{
		SortOptions: sortOptions,
		PriceRanges: priceRanges,
		TimeSlots:   timeSlots,
		SeatTypes:   seatTypes,
		Amenities:   amenities,
		Cities:      cities,
	}, nil
}

func (s *ConstantsServiceImpl) GetCities(ctx context.Context) ([]string, error) {
	return []string{
		"Hà Nội",
		"TP. Hồ Chí Minh",
		"Đà Nẵng",
		"Đà Lạt",
		"Hải Phòng",
		"Cần Thơ",
		"Huế",
		"Nha Trang",
		"Quảng Ninh",
		"Bắc Ninh",
		"Bình Dương",
		"Bình Định",
		"Bình Phước",
		"Bến Tre",
		"Gia Lai",
		"Kiên Giang",
		"Lâm Đồng",
		"Long An",
		"Nam Định",
		"Nghệ An",
		"Phú Yên",
		"Quảng Bình",
		"Quảng Nam",
		"Quảng Ngãi",
		"Sóc Trăng",
		"Tây Ninh",
		"Thái Nguyên",
		"Thanh Hóa",
		"Tiền Giang",
		"Trà Vinh",
		"Vĩnh Long",
		"Vũng Tàu",
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

	searchFilters, err := s.GetSearchFilterConstants(ctx)
	if err != nil {
		return nil, err
	}

	return &model.ConstantsResponse{
		Bus:           *busConstants,
		Route:         *routeConstants,
		Trip:          *tripConstants,
		SearchFilters: *searchFilters,
	}, nil
}

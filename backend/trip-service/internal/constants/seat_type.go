package constants

type SeatType string

const (
	SeatTypeStandard SeatType = "standard"
	SeatTypeVIP      SeatType = "vip"
	SeatTypeSleeper  SeatType = "sleeper"
)

func (s SeatType) String() string {
	return string(s)
}

func (s SeatType) IsValid() bool {
	switch s {
	case SeatTypeStandard, SeatTypeVIP, SeatTypeSleeper:
		return true
	}
	return false
}

// GetPriceMultiplier returns the price multiplier for each seat type
func (s SeatType) GetPriceMultiplier() float64 {
	switch s {
	case SeatTypeStandard:
		return 1.0
	case SeatTypeVIP:
		return 1.5
	case SeatTypeSleeper:
		return 2.0
	default:
		return 1.0
	}
}

// AllSeatTypes returns all valid seat types
func AllSeatTypes() []SeatType {
	return []SeatType{
		SeatTypeStandard,
		SeatTypeVIP,
		SeatTypeSleeper,
	}
}

// GetDisplayName returns a user-friendly display name for the seat type
func (s SeatType) GetDisplayName() string {
	switch s {
	case SeatTypeStandard:
		return "Ghế ngồi thường"
	case SeatTypeVIP:
		return "Ghế ngồi VIP"
	case SeatTypeSleeper:
		return "Giường nằm"
	default:
		return string(s)
	}
}

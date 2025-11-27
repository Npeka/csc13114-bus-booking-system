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

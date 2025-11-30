package constants

type BusType string

const (
	BusTypeStandard     BusType = "standard"
	BusTypeVIP          BusType = "vip"
	BusTypeSleeper      BusType = "sleeper"
	BusTypeDoubleDecker BusType = "double_decker"
)

func (b BusType) String() string {
	return string(b)
}

func (b BusType) IsValid() bool {
	switch b {
	case BusTypeStandard, BusTypeVIP, BusTypeSleeper, BusTypeDoubleDecker:
		return true
	}
	return false
}

// AllBusTypes returns all valid bus types
func AllBusTypes() []BusType {
	return []BusType{
		BusTypeStandard,
		BusTypeVIP,
		BusTypeSleeper,
		BusTypeDoubleDecker,
	}
}

// GetDisplayName returns a user-friendly display name for the bus type
func (b BusType) GetDisplayName() string {
	switch b {
	case BusTypeStandard:
		return "Ghế ngồi thường"
	case BusTypeVIP:
		return "Ghế ngồi VIP"
	case BusTypeSleeper:
		return "Giường nằm"
	case BusTypeDoubleDecker:
		return "Xe 2 tầng"
	default:
		return string(b)
	}
}

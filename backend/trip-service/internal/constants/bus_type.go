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

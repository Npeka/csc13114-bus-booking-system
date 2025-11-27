package constants

type StopType string

const (
	StopTypePickup  StopType = "pickup"
	StopTypeDropoff StopType = "dropoff"
	StopTypeBoth    StopType = "both"
)

func (s StopType) String() string {
	return string(s)
}

func (s StopType) IsValid() bool {
	switch s {
	case StopTypePickup, StopTypeDropoff, StopTypeBoth:
		return true
	}
	return false
}

// AllStopTypes returns all valid stop types
func AllStopTypes() []StopType {
	return []StopType{
		StopTypePickup,
		StopTypeDropoff,
		StopTypeBoth,
	}
}

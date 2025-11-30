package constants

type ConstantType string

const (
	ConstantTypeBus           ConstantType = "bus"
	ConstantTypeRoute         ConstantType = "route"
	ConstantTypeTrip          ConstantType = "trip"
	ConstantTypeSearchFilters ConstantType = "search_filters"
	ConstantTypeCities        ConstantType = "cities"
)

func (c ConstantType) String() string {
	return string(c)
}

func (c ConstantType) IsValid() bool {
	switch c {
	case ConstantTypeBus, ConstantTypeRoute, ConstantTypeTrip, ConstantTypeSearchFilters, ConstantTypeCities:
		return true
	}
	return false
}

// AllConstantTypes returns all valid constant types
func AllConstantTypes() []ConstantType {
	return []ConstantType{
		ConstantTypeBus,
		ConstantTypeRoute,
		ConstantTypeTrip,
		ConstantTypeSearchFilters,
		ConstantTypeCities,
	}
}

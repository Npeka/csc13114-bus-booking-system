package constants

type Amenity string

const (
	AmenityWiFi     Amenity = "wifi"
	AmenityAC       Amenity = "ac"
	AmenityToilet   Amenity = "toilet"
	AmenityTV       Amenity = "tv"
	AmenityCharging Amenity = "charging"
	AmenityBlanket  Amenity = "blanket"
	AmenityWater    Amenity = "water"
	AmenitySnacks   Amenity = "snacks"
)

// AllAmenities returns all available amenities
var AllAmenities = []Amenity{
	AmenityWiFi,
	AmenityAC,
	AmenityToilet,
	AmenityTV,
	AmenityCharging,
	AmenityBlanket,
	AmenityWater,
	AmenitySnacks,
}

func (a Amenity) String() string {
	return string(a)
}

func (a Amenity) IsValid() bool {
	for _, amenity := range AllAmenities {
		if a == amenity {
			return true
		}
	}
	return false
}

// GetDisplayName returns a user-friendly display name for the amenity
func (a Amenity) GetDisplayName() string {
	switch a {
	case AmenityWiFi:
		return "Wi-Fi"
	case AmenityAC:
		return "Air Conditioning"
	case AmenityToilet:
		return "Toilet"
	case AmenityTV:
		return "TV/Entertainment"
	case AmenityCharging:
		return "Charging Ports"
	case AmenityBlanket:
		return "Blanket & Pillow"
	case AmenityWater:
		return "Complimentary Water"
	case AmenitySnacks:
		return "Snacks"
	default:
		return string(a)
	}
}

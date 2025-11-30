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
		return "Điều hòa"
	case AmenityToilet:
		return "Nhà vệ sinh"
	case AmenityTV:
		return "TV/Giải trí"
	case AmenityCharging:
		return "Cổng sạc điện thoại"
	case AmenityBlanket:
		return "Chăn & Gối"
	case AmenityWater:
		return "Nước uống miễn phí"
	case AmenitySnacks:
		return "Đồ ăn nhẹ"
	default:
		return string(a)
	}
}

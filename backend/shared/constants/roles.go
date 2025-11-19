package constants

// UserRole represents user roles using bit flags
type UserRole int

const (
	RolePassenger UserRole = 1 << iota // bit 0: 1
	RoleAdmin                          // bit 1: 2
	RoleOperator                       // bit 2: 4
	RoleSupport                        // bit 3: 8
)

// Role constants as int for easier usage
const (
	RolePassengerInt = int(RolePassenger) // 1
	RoleAdminInt     = int(RoleAdmin)     // 2
	RoleOperatorInt  = int(RoleOperator)  // 4
	RoleSupportInt   = int(RoleSupport)   // 8
)

// HasRole checks if a role has a specific permission
func (r UserRole) HasRole(role UserRole) bool {
	return r&role != 0
}

// String returns the string representation of the role
func (r UserRole) String() string {
	switch r {
	case RolePassenger:
		return "passenger"
	case RoleAdmin:
		return "admin"
	case RoleOperator:
		return "operator"
	case RoleSupport:
		return "support"
	default:
		return "unknown"
	}
}

// ToInt converts UserRole to int
func (r UserRole) ToInt() int {
	return int(r)
}

// FromInt creates UserRole from int
func FromInt(role int) UserRole {
	return UserRole(role)
}

// ValidateRole checks if the role value is valid
func ValidateRole(role int) bool {
	return role == RolePassengerInt || role == RoleAdminInt ||
		role == RoleOperatorInt || role == RoleSupportInt
}

package constants

// UserRole represents user roles using bit flags
type UserRole int

const (
	RoleGuest     UserRole = 1 << iota // bit 0: 1
	RolePassenger                      // bit 1: 2
	RoleAdmin                          // bit 2: 4
)

// Role constants as int for easier usage
const (
	RoleGuestInt     = int(RoleGuest)     // 1
	RolePassengerInt = int(RolePassenger) // 2
	RoleAdminInt     = int(RoleAdmin)     // 4
)

// HasRole checks if a role has a specific permission
func (r UserRole) HasRole(role UserRole) bool {
	return r&role != 0
}

// String returns the string representation of the role
func (r UserRole) String() string {
	switch r {
	case RoleGuest:
		return "guest"
	case RolePassenger:
		return "passenger"
	case RoleAdmin:
		return "admin"
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

// UserStatus represents user account status
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"    // User account is active and can login
	UserStatusInactive  UserStatus = "inactive"  // User account is inactive (not yet activated)
	UserStatusSuspended UserStatus = "suspended" // User account is suspended (temporarily blocked)
	UserStatusVerified  UserStatus = "verified"  // User account is verified (via Firebase/Email)
)

// IsValid checks if the status is valid
func (s UserStatus) IsValid() bool {
	switch s {
	case UserStatusActive, UserStatusInactive, UserStatusSuspended, UserStatusVerified:
		return true
	default:
		return false
	}
}

// String returns the string representation of the status
func (s UserStatus) String() string {
	return string(s)
}

// CanLogin checks if user with this status can login
func (s UserStatus) CanLogin() bool {
	return s == UserStatusActive || s == UserStatusVerified
}

// ValidateRole checks if the role value is valid
func ValidateRole(role int) bool {
	return role == RoleGuestInt || role == RolePassengerInt || role == RoleAdminInt
}

// FromString converts role string to UserRole
func FromString(role string) UserRole {
	switch role {
	case "guest":
		return RoleGuest
	case "passenger":
		return RolePassenger
	case "admin":
		return RoleAdmin
	default:
		return 0
	}
}

// FromStringSlice converts slice of role strings to slice of UserRole
func FromStringSlice(roles []string) []UserRole {
	result := make([]UserRole, 0, len(roles))
	for _, role := range roles {
		if r := FromString(role); r != 0 {
			result = append(result, r)
		}
	}
	return result
}

// ToStringSlice converts slice of UserRole to slice of strings
func ToStringSlice(roles []UserRole) []string {
	result := make([]string, len(roles))
	for i, role := range roles {
		result[i] = role.String()
	}
	return result
}

// IsValidRoleString checks if a role string is valid
func IsValidRoleString(role string) bool {
	return FromString(role) != 0
}

// ValidRoleStrings returns all valid role strings
func ValidRoleStrings() []string {
	return []string{"guest", "passenger", "admin"}
}

// HasAnyRole checks if a role has any of the specified roles
func (r UserRole) HasAnyRole(roles []UserRole) bool {
	for _, role := range roles {
		if r.HasRole(role) {
			return true
		}
	}
	return false
}

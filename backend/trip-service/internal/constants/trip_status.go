package constants

type TripStatus string

const (
	TripStatusScheduled  TripStatus = "scheduled"
	TripStatusInProgress TripStatus = "in_progress"
	TripStatusCompleted  TripStatus = "completed"
	TripStatusCancelled  TripStatus = "cancelled"
	TripStatusDelayed    TripStatus = "delayed"
)

func (t TripStatus) String() string {
	return string(t)
}

func (t TripStatus) IsValid() bool {
	switch t {
	case TripStatusScheduled, TripStatusInProgress, TripStatusCompleted,
		TripStatusCancelled, TripStatusDelayed:
		return true
	}
	return false
}

// AllTripStatuses returns all valid trip statuses
func AllTripStatuses() []TripStatus {
	return []TripStatus{
		TripStatusScheduled,
		TripStatusInProgress,
		TripStatusCompleted,
		TripStatusCancelled,
		TripStatusDelayed,
	}
}

// GetDisplayName returns a user-friendly display name for the trip status
func (t TripStatus) GetDisplayName() string {
	switch t {
	case TripStatusScheduled:
		return "Đã lên lịch"
	case TripStatusInProgress:
		return "Đang di chuyển"
	case TripStatusCompleted:
		return "Hoàn thành"
	case TripStatusCancelled:
		return "Đã hủy"
	case TripStatusDelayed:
		return "Chậm trễ"
	default:
		return string(t)
	}
}

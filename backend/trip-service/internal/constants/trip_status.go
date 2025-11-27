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

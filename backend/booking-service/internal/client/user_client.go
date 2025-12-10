package client

import (
	"bus-booking/booking-service/internal/model/user"
	"bus-booking/shared/client"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UserClient interface for user service communication
type UserClient interface {
	CreateGuest(ctx context.Context, req *user.CreateGuestRequest) (*user.GuestResponse, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*user.User, error)
}

type userClientImpl struct {
	http client.HTTPClient
}

func NewUserClient(serviceName, baseURL string) UserClient {
	httpClient := client.NewHTTPClient(&client.Config{
		ServiceName: serviceName,
		BaseURL:     baseURL,
	})

	return &userClientImpl{
		http: httpClient,
	}
}

// Request/Response types

// CreateGuest creates a guest user account
func (c *userClientImpl) CreateGuest(ctx context.Context, req *user.CreateGuestRequest) (*user.GuestResponse, error) {
	endpoint := "/api/v1/auth/guest"

	res, err := c.http.Post(ctx, endpoint, req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	guestData, err := client.ParseData[user.GuestResponse](res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trip response: %w", err)
	}

	return guestData, nil
}

func (c *userClientImpl) GetUser(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	// Use internal endpoint for service-to-service calls
	endpoint := fmt.Sprintf("/api/v1/internal/users/%s", userID.String())

	res, err := c.http.Get(ctx, endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	userData, err := client.ParseData[user.User](res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	return userData, nil
}

package middleware

import (
	"context"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"bus-booking/shared/constants"
	"bus-booking/user-service/internal/model"
	"bus-booking/user-service/internal/repository"
)

type FirebaseAuthMiddleware struct {
	firebaseAuth *auth.Client
	userRepo     repository.UserRepository
}

func NewFirebaseAuthMiddleware(
	firebaseAuth *auth.Client,
	userRepo repository.UserRepository,
) *FirebaseAuthMiddleware {
	return &FirebaseAuthMiddleware{
		firebaseAuth: firebaseAuth,
		userRepo:     userRepo,
	}
}

func (m *FirebaseAuthMiddleware) FirebaseAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.firebaseAuth == nil {
			log.Warn().Msg("Firebase Auth client not initialized - skipping Firebase auth")
			c.Next()
			return
		}

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		idToken := parts[1]

		// Verify the Firebase ID token
		token, err := m.firebaseAuth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			log.Warn().Err(err).Msg("Invalid Firebase ID token")
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Try to get user from database using Firebase UID
		user, err := m.userRepo.GetByFirebaseUID(context.Background(), token.UID)
		if err != nil {
			log.Error().Err(err).Str("firebase_uid", token.UID).Msg("Failed to get user by Firebase UID")
			c.JSON(500, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}

		// If user doesn't exist, create a new one based on Firebase token
		if user == nil {
			email, ok := token.Claims["email"].(string)
			if !ok || email == "" {
				log.Error().Str("firebase_uid", token.UID).Msg("No email in Firebase token")
				c.JSON(400, gin.H{"error": "Email is required"})
				c.Abort()
				return
			}

			name, _ := token.Claims["name"].(string)
			if name == "" {
				name = email
			}

			picture, _ := token.Claims["picture"].(string)
			emailVerified, _ := token.Claims["email_verified"].(bool)

			user = &model.User{
				Email:         email,
				FullName:      name,
				Avatar:        picture,
				Role:          constants.RolePassenger,
				Status:        "verified",
				FirebaseUID:   token.UID,
				EmailVerified: emailVerified,
				PhoneVerified: false,
			}

			if err := m.userRepo.Create(context.Background(), user); err != nil {
				log.Error().Err(err).Str("firebase_uid", token.UID).Msg("Failed to create user from Firebase token")
				c.JSON(500, gin.H{"error": "Authentication failed"})
				c.Abort()
				return
			}
		}

		// Check if user is active
		if user.Status != "active" && user.Status != "verified" {
			c.JSON(401, gin.H{"error": "Account is not active"})
			c.Abort()
			return
		}

		// Store user info and Firebase token in context
		c.Set("user", user)
		c.Set("firebase_uid", token.UID)
		c.Set("firebase_token", token)

		log.Debug().
			Str("user_id", user.ID.String()).
			Str("firebase_uid", token.UID).
			Str("email", user.Email).
			Msg("Firebase authentication successful")

		c.Next()
	}
}

func (m *FirebaseAuthMiddleware) OptionalFirebaseAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.firebaseAuth == nil {
			log.Debug().Msg("Firebase Auth client not initialized - skipping optional Firebase auth")
			c.Next()
			return
		}

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		idToken := parts[1]

		// Verify the Firebase ID token
		token, err := m.firebaseAuth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			log.Debug().Err(err).Msg("Invalid Firebase ID token in optional auth")
			c.Next()
			return
		}

		// Try to get user from database using Firebase UID
		user, err := m.userRepo.GetByFirebaseUID(context.Background(), token.UID)
		if err != nil {
			log.Debug().Err(err).Str("firebase_uid", token.UID).Msg("Failed to get user by Firebase UID in optional auth")
			c.Next()
			return
		}

		if user != nil && (user.Status == "active" || user.Status == "verified") {
			// Store user info and Firebase token in context
			c.Set("user", user)
			c.Set("firebase_uid", token.UID)
			c.Set("firebase_token", token)

			log.Debug().
				Str("user_id", user.ID.String()).
				Str("firebase_uid", token.UID).
				Str("email", user.Email).
				Msg("Optional Firebase authentication successful")
		}

		c.Next()
	}
}

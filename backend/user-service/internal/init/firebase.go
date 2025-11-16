package appinit

import (
	"context"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"

	"bus-booking/user-service/config"
)

func InitFirebase(cfg *config.Config) (*auth.Client, error) {
	ctx := context.Background()
	var firebaseApp *firebase.App
	var err error

	// Try to load from service account file first
	serviceAccountPath := cfg.Firebase.ServiceAccountKeyPath

	// Check if service account file exists at relative path
	if _, err := os.Stat(serviceAccountPath); os.IsNotExist(err) {
		// Try absolute path from project root
		projectRoot := getProjectRoot()
		absolutePath := filepath.Join(projectRoot, serviceAccountPath)
		if _, err := os.Stat(absolutePath); err == nil {
			serviceAccountPath = absolutePath
		}
	}

	// Try loading from JSON file first
	if _, err := os.Stat(serviceAccountPath); err == nil {
		opt := option.WithCredentialsFile(serviceAccountPath)
		firebaseConfig := &firebase.Config{
			ProjectID: cfg.Firebase.ProjectID,
		}

		firebaseApp, err = firebase.NewApp(ctx, firebaseConfig, opt)
		if err == nil {
			log.Info().
				Str("service_account_path", serviceAccountPath).
				Str("project_id", firebaseConfig.ProjectID).
				Msg("Firebase initialized with service account JSON file")
		} else {
			log.Warn().Err(err).
				Str("service_account_path", serviceAccountPath).
				Msg("Failed to initialize Firebase with JSON file, trying environment variables")
		}
	}

	// Fallback to environment variables if JSON file failed or doesn't exist
	if firebaseApp == nil {
		firebaseApp, err = firebase.NewApp(ctx, nil)
		if err != nil {
			log.Warn().Err(err).Msg("Firebase initialization failed completely - running without Firebase Auth")
			return nil, nil
		}
		log.Info().Msg("Firebase initialized with environment variables")
	}

	authClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("Firebase Auth client initialization failed - running without Firebase Auth")
		return nil, nil
	}

	log.Info().Msg("Firebase Auth client initialized successfully")
	return authClient, nil
}

// getProjectRoot returns the project root directory
func getProjectRoot() string {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Navigate up to find project root (where go.mod exists)
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return wd
}

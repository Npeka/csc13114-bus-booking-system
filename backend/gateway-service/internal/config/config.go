package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig             `yaml:"server"`
	Services map[string]ServiceConfig `yaml:"services"`
	Auth     AuthConfig               `yaml:"auth"`
	CORS     CORSConfig               `yaml:"cors"`
}

type ServerConfig struct {
	Port         string `yaml:"port" env:"PORT" default:"8000"`
	Host         string `yaml:"host" env:"HOST" default:"0.0.0.0"`
	ReadTimeout  int    `yaml:"read_timeout" default:"30"`
	WriteTimeout int    `yaml:"write_timeout" default:"30"`
}

type ServiceConfig struct {
	URL     string `yaml:"url"`
	Timeout int    `yaml:"timeout" default:"30"`
	Retries int    `yaml:"retries" default:"3"`
}

type AuthConfig struct {
	UserServiceURL string `yaml:"user_service_url" env:"USER_SERVICE_URL"`
	VerifyEndpoint string `yaml:"verify_endpoint" default:"/api/v1/auth/verify-token"`
	Timeout        int    `yaml:"timeout" default:"5"`
}

type CORSConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

type RouteConfig struct {
	Service string  `yaml:"service"`
	Routes  []Route `yaml:"routes"`
}

type Route struct {
	Path        string            `yaml:"path"`
	Methods     []string          `yaml:"methods"`
	Service     string            `yaml:"service"`
	Auth        *AuthRequirement  `yaml:"auth,omitempty"`
	Headers     map[string]string `yaml:"headers,omitempty"`
	StripPrefix string            `yaml:"strip_prefix,omitempty"`
	Rewrite     *RewriteRule      `yaml:"rewrite,omitempty"`
}

type AuthRequirement struct {
	Required bool     `yaml:"required"`
	Roles    []string `yaml:"roles,omitempty"`
}

type RewriteRule struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

// LoadConfig loads configuration from file and environment
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnvOrDefault("PORT", "8000"),
			Host:         getEnvOrDefault("HOST", "0.0.0.0"),
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Auth: AuthConfig{
			UserServiceURL: getEnvOrDefault("USER_INTERNAL_URL", "http://user-internal.cluster.local"),
			VerifyEndpoint: "/api/v1/auth/verify-token",
			Timeout:        5,
		},
		CORS: CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           86400,
		},
		Services: map[string]ServiceConfig{
			"user": {
				URL:     getEnvOrDefault("USER_INTERNAL_URL", "http://user-internal.cluster.local"),
				Timeout: 30,
				Retries: 3,
			},
			"trips": {
				URL:     getEnvOrDefault("TRIPS_INTERNAL_URL", "http://trips-internal.cluster.local"),
				Timeout: 30,
				Retries: 3,
			},
			"bookings": {
				URL:     getEnvOrDefault("BOOKINGS_INTERNAL_URL", "http://bookings-internal.cluster.local"),
				Timeout: 30,
				Retries: 3,
			},
			"templates": {
				URL:     getEnvOrDefault("TEMPLATES_INTERNAL_URL", "http://templates-internal.cluster.local"),
				Timeout: 30,
				Retries: 3,
			},
			"payments": {
				URL:     getEnvOrDefault("PAYMENTS_INTERNAL_URL", "http://payments-internal.cluster.local"),
				Timeout: 30,
				Retries: 3,
			},
		},
	}

	// Load from config file if exists
	if configPath != "" {
		if err := loadFromFile(config, configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	return config, nil
}

// LoadRoutes loads route configurations from directory
func LoadRoutes(routesDir string) (*RouteConfig, error) {
	routes := &RouteConfig{Routes: []Route{}}

	err := filepath.Walk(routesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			routeConfig := &RouteConfig{}
			if err := loadFromFile(routeConfig, path); err != nil {
				return fmt.Errorf("failed to load route config from %s: %w", path, err)
			}

			// Set service name for all routes in this file
			for i := range routeConfig.Routes {
				if routeConfig.Routes[i].Service == "" && routeConfig.Service != "" {
					routeConfig.Routes[i].Service = routeConfig.Service
				}
			}

			routes.Routes = append(routes.Routes, routeConfig.Routes...)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load routes: %w", err)
	}

	return routes, nil
}

func loadFromFile(config interface{}, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

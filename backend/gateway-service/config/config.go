package config

import (
	"fmt"
	"os"
	"path/filepath"

	sharedConfig "bus-booking/shared/config"
	"bus-booking/shared/constants"

	"gopkg.in/yaml.v3"
)

type Config struct {
	*sharedConfig.BaseConfig
	Services    ServicesConfig `envPrefix:"SERVICES_"`
	Auth        AuthConfig     `envPrefix:"AUTH_"`
	ServicesMap map[string]ServiceConfig
}

type ServicesConfig struct {
	User    ServiceConfig `envPrefix:"USER_"`
	Trip    ServiceConfig `envPrefix:"TRIP_"`
	Booking ServiceConfig `envPrefix:"BOOKING_"`
	Payment ServiceConfig `envPrefix:"PAYMENT_"`
}

type ServiceConfig struct {
	URL     string `env:"URL"`
	Timeout int    `env:"TIMEOUT" envDefault:"30"`
	Retries int    `env:"RETRIES" envDefault:"3"`
}

type AuthConfig struct {
	UserServiceURL string `env:"USER_SERVICE_URL" envDefault:"http://localhost:8080"`
	VerifyEndpoint string `env:"VERIFY_ENDPOINT" envDefault:"/api/v1/auth/verify-token"`
	Timeout        int    `env:"TIMEOUT" envDefault:"60"`
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

func (c *Config) BuildServiceMap() map[string]ServiceConfig {
	return map[string]ServiceConfig{
		"user":    c.Services.User,
		"trip":    c.Services.Trip,
		"booking": c.Services.Booking,
		"payment": c.Services.Payment,
	}
}

func MustLoadConfig() *Config {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	return config
}

// LoadConfig loads configuration from environment variables and file
func LoadConfig() (*Config, error) {
	// Load config using shared pattern
	config, err := sharedConfig.LoadConfig[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Build services map for easy lookup
	config.ServicesMap = config.BuildServiceMap()

	return config, nil
}

func MustLoadRoutes() *RouteConfig {
	routes, err := LoadRoutes()
	if err != nil {
		panic(err)
	}
	return routes
}

// LoadRoutes loads route configurations from directory
func LoadRoutes() (*RouteConfig, error) {
	routes := &RouteConfig{Routes: []Route{}}

	err := filepath.Walk("routes", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			routeConfig := &RouteConfig{}
			if err := loadFromFile(routeConfig, path); err != nil {
				return fmt.Errorf("failed to load route config from %s: %w", path, err)
			}

			// Set service name for all routes in this file and validate roles
			for i := range routeConfig.Routes {
				if routeConfig.Routes[i].Service == "" && routeConfig.Service != "" {
					routeConfig.Routes[i].Service = routeConfig.Service
				}

				// Validate roles if auth is configured
				if routeConfig.Routes[i].Auth != nil {
					for _, role := range routeConfig.Routes[i].Auth.Roles {
						if !constants.IsValidRoleString(role) {
							return fmt.Errorf("invalid role '%s' in route %s:%v", role, routeConfig.Routes[i].Path, routeConfig.Routes[i].Methods)
						}
					}
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
	// #nosec G304 -- path is from trusted routes directory walked by filepath.Walk
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

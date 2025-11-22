package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"bus-booking/gateway-service/config"
	"bus-booking/gateway-service/internal/auth"
)

type Gateway struct {
	config     *config.Config
	routes     *config.RouteConfig
	authClient *auth.Client
	httpClient *http.Client
}

func NewGateway(cfg *config.Config, routes *config.RouteConfig) *Gateway {
	return &Gateway{
		config:     cfg,
		routes:     routes,
		authClient: auth.NewClient(&cfg.Auth),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (g *Gateway) SetupRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "gateway-service",
		})
	})

	// Setup proxy routes
	log.Printf("Total routes to setup: %d", len(g.routes.Routes))
	for i, route := range g.routes.Routes {
		// Create route with service prefix
		prefixedPath := "/" + route.Service + route.Path
		log.Printf("Route %d: Setting up route: %s -> %s %v -> service: %s", i+1, route.Path, prefixedPath, route.Methods, route.Service)

		// Create new route with prefixed path
		prefixedRoute := route
		prefixedRoute.Path = prefixedPath

		// Validate route path before setting up
		if err := g.validateRoutePath(prefixedRoute.Path); err != nil {
			log.Printf("Skipping invalid route %s: %v", prefixedRoute.Path, err)
			continue
		}

		// Add extra safety check
		if strings.HasSuffix(prefixedRoute.Path, "/*") {
			log.Printf("ERROR: Route %s ends with /* which is invalid for Gin router", prefixedRoute.Path)
			continue
		}

		g.setupRoute(router, prefixedRoute)
		log.Printf("Successfully registered route %d: %s", i+1, prefixedRoute.Path)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "route not found",
			"path":    c.Request.URL.Path,
			"method":  c.Request.Method,
			"message": "The requested route is not configured in the gateway",
		})
	})
}

func (g *Gateway) setupRoute(router *gin.Engine, route config.Route) {
	handler := g.createProxyHandler(route)

	for _, method := range route.Methods {
		methodUpper := strings.ToUpper(method)
		log.Printf("Registering %s %s", methodUpper, route.Path)

		switch methodUpper {
		case "GET":
			router.GET(route.Path, handler)
		case "POST":
			router.POST(route.Path, handler)
		case "PUT":
			router.PUT(route.Path, handler)
		case "DELETE":
			router.DELETE(route.Path, handler)
		case "PATCH":
			router.PATCH(route.Path, handler)
		case "OPTIONS":
			router.OPTIONS(route.Path, handler)
		case "HEAD":
			router.HEAD(route.Path, handler)
		}
	}
}

func (g *Gateway) createProxyHandler(route config.Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check authentication if required
		var userContext *auth.UserContext
		if route.Auth != nil && route.Auth.Required {
			var err error
			userContext, err = g.authenticateRequest(c)
			if err != nil {
				c.JSON(401, gin.H{
					"error":   "authentication failed",
					"message": err.Error(),
				})
				return
			}

			// Check roles if specified
			if len(route.Auth.Roles) > 0 && !userContext.HasAnyRoleString(route.Auth.Roles) {
				c.JSON(403, gin.H{
					"error":   "insufficient permissions",
					"message": fmt.Sprintf("required roles: %v", route.Auth.Roles),
				})
				return
			}
		}

		// Get service configuration (case-insensitive lookup)
		serviceConfig, exists := g.getServiceConfig(route.Service)
		log.Printf("Service config for %s: %+v", route.Service, serviceConfig)
		if !exists {
			c.JSON(500, gin.H{
				"error":   "service not configured",
				"service": route.Service,
			})
			return
		}

		// Build target URL
		targetURL, err := g.buildTargetURL(serviceConfig, route, c)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   "failed to build target URL",
				"message": err.Error(),
			})
			return
		}

		// Proxy the request
		if err := g.proxyRequest(c, targetURL, userContext, route); err != nil {
			c.JSON(502, gin.H{
				"error":   "proxy request failed",
				"message": err.Error(),
			})
			return
		}
	}
}

func (g *Gateway) authenticateRequest(c *gin.Context) (*auth.UserContext, error) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("missing authorization header")
	}

	// Verify token with user service
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(g.config.Auth.Timeout)*time.Second)
	defer cancel()

	return g.authClient.VerifyToken(ctx, authHeader)
}

func (g *Gateway) buildTargetURL(serviceConfig config.ServiceConfig, route config.Route, c *gin.Context) (string, error) {
	// Start with service base URL
	baseURL := serviceConfig.URL

	// Get actual request path
	actualPath := c.Request.URL.Path

	// Strip service prefix from the actual path to get target path
	// If route.Path starts with /service/, strip it to get the original API path
	servicePrefix := "/" + route.Service
	targetPath := actualPath

	if strings.HasPrefix(actualPath, servicePrefix) {
		// Strip service prefix: /user/api/v1/auth -> /api/v1/auth
		targetPath = strings.TrimPrefix(actualPath, servicePrefix)
		log.Printf("Stripped service prefix: %s -> %s", actualPath, targetPath)
	}

	// Handle path rewriting
	if route.Rewrite != nil {
		re, err := regexp.Compile(route.Rewrite.From)
		if err != nil {
			return "", fmt.Errorf("invalid rewrite regex: %w", err)
		}
		targetPath = re.ReplaceAllString(targetPath, route.Rewrite.To)
	}

	// Handle additional prefix stripping
	if route.StripPrefix != "" {
		targetPath = strings.TrimPrefix(targetPath, route.StripPrefix)
	}

	// Ensure target path starts with /
	if !strings.HasPrefix(targetPath, "/") {
		targetPath = "/" + targetPath
	}

	// Build full URL
	fullURL := strings.TrimSuffix(baseURL, "/") + targetPath

	// Add query parameters
	if c.Request.URL.RawQuery != "" {
		fullURL += "?" + c.Request.URL.RawQuery
	}

	return fullURL, nil
}

func (g *Gateway) proxyRequest(c *gin.Context, targetURL string, userContext *auth.UserContext, route config.Route) error {
	// Create new request
	req, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Copy headers from original request
	for key, values := range c.Request.Header {
		// Skip some headers that shouldn't be proxied
		if g.shouldSkipHeader(key) {
			continue
		}
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Add user context headers if authenticated
	if userContext != nil {
		userHeaders := userContext.ToHeaders()
		for key, value := range userHeaders {
			req.Header.Set(key, value)
		}
	}

	// Add custom headers from route config
	for key, value := range route.Headers {
		req.Header.Set(key, value)
	}

	// Set gateway identification header
	req.Header.Set("X-Gateway", "bus-booking-gateway")
	req.Header.Set("X-Forwarded-For", c.ClientIP())
	req.Header.Set("X-Forwarded-Host", c.Request.Host)
	req.Header.Set("X-Forwarded-Proto", g.getScheme(c))

	// Make the request
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to proxy request: %w", err)
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Set status code
	c.Status(resp.StatusCode)

	// Copy response body
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy response body: %w", err)
	}

	return nil
}

func (g *Gateway) shouldSkipHeader(header string) bool {
	skipHeaders := []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}

	headerLower := strings.ToLower(header)
	for _, skip := range skipHeaders {
		if strings.ToLower(skip) == headerLower {
			return true
		}
	}
	return false
}

func (g *Gateway) getScheme(c *gin.Context) string {
	if c.Request.TLS != nil {
		return "https"
	}
	if scheme := c.GetHeader("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	return "http"
}

// validateRoutePath validates that the route path is properly formatted for Gin router
func (g *Gateway) validateRoutePath(path string) error {
	// Check for unnamed wildcards
	if strings.Contains(path, "/*") {
		// Find all occurrences of /*
		for i := 0; i < len(path)-1; i++ {
			if path[i] == '/' && path[i+1] == '*' {
				// Check if there's a name after the *
				if i+2 >= len(path) {
					return fmt.Errorf("wildcard must be named (use /*filepath instead of /*)")
				}
				// If the next character is not alphanumeric, it's an unnamed wildcard
				if i+2 < len(path) && !isAlphaNumeric(path[i+2]) && path[i+2] != '_' {
					return fmt.Errorf("wildcard must be named (use /*filepath instead of /*)")
				}
			}
		}
	}
	return nil
}

func isAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

// getServiceConfig performs case-insensitive lookup for service configuration
func (g *Gateway) getServiceConfig(serviceName string) (config.ServiceConfig, bool) {
	// Try exact match first
	if serviceConfig, exists := g.config.ServicesMap[serviceName]; exists {
		return serviceConfig, true
	}

	// Try case-insensitive match
	for configKey, serviceConfig := range g.config.ServicesMap {
		if strings.EqualFold(configKey, serviceName) {
			return serviceConfig, true
		}
	}

	return config.ServiceConfig{}, false
}

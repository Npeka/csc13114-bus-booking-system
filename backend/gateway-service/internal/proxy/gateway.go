package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"bus-booking/gateway-service/internal/auth"
	"bus-booking/gateway-service/internal/config"
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

// SetupRoutes configures all gateway routes
func (g *Gateway) SetupRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "gateway-service",
		})
	})

	// Setup proxy routes
	for _, route := range g.routes.Routes {
		g.setupRoute(router, route)
	}

	// Catch-all route for unmatched paths
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

	// Register route for each method
	for _, method := range route.Methods {
		switch strings.ToUpper(method) {
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
			if len(route.Auth.Roles) > 0 && !userContext.HasAnyRole(route.Auth.Roles) {
				c.JSON(403, gin.H{
					"error":   "insufficient permissions",
					"message": fmt.Sprintf("required roles: %v", route.Auth.Roles),
				})
				return
			}
		}

		// Get service configuration
		serviceConfig, exists := g.config.Services[route.Service]
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

	// Use route path as target path (since target field is removed)
	targetPath := route.Path

	// Handle path rewriting
	if route.Rewrite != nil {
		re, err := regexp.Compile(route.Rewrite.From)
		if err != nil {
			return "", fmt.Errorf("invalid rewrite regex: %w", err)
		}
		targetPath = re.ReplaceAllString(targetPath, route.Rewrite.To)
	}

	// Handle prefix stripping
	if route.StripPrefix != "" {
		targetPath = strings.TrimPrefix(targetPath, route.StripPrefix)
	}

	// For parameterized routes, replace with actual values
	if strings.Contains(targetPath, ":") {
		// Extract params from current request
		actualPath := c.Request.URL.Path
		targetPath = actualPath
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

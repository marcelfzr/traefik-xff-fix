// Package traefik_xff_fix a Traefik plugin that rewrites X-Forwarded-For to contain only the leftmost client IP without port.
package traefik_xff_fix

import (
	"context"
	"net"
	"net/http"
	"strings"
)

const xForwardedFor = "X-Forwarded-For"

// Config the plugin configuration.
type Config struct {
	// Empty config - plugin has fixed behavior.
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// XFFFix a Traefik middleware that fixes the X-Forwarded-For header.
type XFFFix struct {
	next http.Handler
	name string
}

// New created a new XFFFix plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &XFFFix{
		next: next,
		name: name,
	}, nil
}

func (x *XFFFix) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	value := req.Header.Get(xForwardedFor)
	if value != "" {
		// Split by comma and take the leftmost (original client) IP
		parts := strings.SplitN(value, ",", 2)
		leftmost := strings.TrimSpace(parts[0])

		// Strip port if present
		host, _, err := net.SplitHostPort(leftmost)
		if err == nil {
			leftmost = host
		}

		req.Header.Set(xForwardedFor, leftmost)
	}

	x.next.ServeHTTP(rw, req)
}

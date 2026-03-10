// Package traefik_xff_fix a Traefik plugin that rewrites X-Forwarded-For to contain only the leftmost client IP without port.
package traefik_xff_fix

import (
	"context"
	"net"
	"net/http"
	"strconv"
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
func New(_ context.Context, next http.Handler, _ *Config, name string) (http.Handler, error) {
	return &XFFFix{
		next: next,
		name: name,
	}, nil
}

func (x *XFFFix) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	value := req.Header.Get(xForwardedFor)
	if value != "" {
		leftmost := normalizedLeftmostIP(value)
		if leftmost != "" {
			// Let Traefik generate X-Forwarded-For from this rewritten remote address
			// so the backend receives only one client IP.
			req.Header.Del(xForwardedFor)
			req.RemoteAddr = net.JoinHostPort(leftmost, strconv.Itoa(0))
		}
	}

	x.next.ServeHTTP(rw, req)
}

func normalizedLeftmostIP(xffValue string) string {
	parts := strings.SplitN(xffValue, ",", 2)

	leftmost := strings.TrimSpace(parts[0])

	if leftmost == "" {
		return ""
	}

	// Strip port when header contains host:port or [ipv6]:port.
	host, _, err := net.SplitHostPort(leftmost)
	if err == nil {
		return host
	}

	// Handle bracketed IPv6 literals without ports.
	return strings.Trim(leftmost, "[]")
}

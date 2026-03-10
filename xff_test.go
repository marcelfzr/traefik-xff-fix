package traefik_xff_fix_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	traefik_xff_fix "github.com/marcelfzr/traefik-xff-fix"
)

func TestXFFFix_MultipleIPs(t *testing.T) {
	handler := mustNewHandler(t)
	req := mustNewRequest(t, "203.0.113.195, 70.41.3.18, 150.172.238.178")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertXFFDeleted(t, req)
	assertRemoteAddrHost(t, req, "203.0.113.195")
}

func TestXFFFix_SingleIPWithPort(t *testing.T) {
	handler := mustNewHandler(t)
	req := mustNewRequest(t, "203.0.113.195:8080")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertXFFDeleted(t, req)
	assertRemoteAddrHost(t, req, "203.0.113.195")
}

func TestXFFFix_IPv6WithBracketPort(t *testing.T) {
	handler := mustNewHandler(t)
	req := mustNewRequest(t, "[2001:db8::1]:8080, 10.0.0.1")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertXFFDeleted(t, req)
	assertRemoteAddrHost(t, req, "2001:db8::1")
}

func TestXFFFix_PlainIPv6(t *testing.T) {
	handler := mustNewHandler(t)
	req := mustNewRequest(t, "2001:db8::1, 10.0.0.1")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertXFFDeleted(t, req)
	assertRemoteAddrHost(t, req, "2001:db8::1")
}

func TestXFFFix_SingleIPv4NoPort(t *testing.T) {
	handler := mustNewHandler(t)
	req := mustNewRequest(t, "192.168.1.1")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertXFFDeleted(t, req)
	assertRemoteAddrHost(t, req, "192.168.1.1")
}

func TestXFFFix_EmptyHeader(t *testing.T) {
	handler := mustNewHandler(t)
	req := mustNewRequest(t, "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Empty/missing header: no crash, passthrough.
	assertXFF(t, req, "")
}

func TestXFFFix_SpacesAroundCommas(t *testing.T) {
	handler := mustNewHandler(t)
	req := mustNewRequest(t, "  10.0.0.1  ,  10.0.0.2  ")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertXFFDeleted(t, req)
	assertRemoteAddrHost(t, req, "10.0.0.1")
}

func mustNewHandler(t *testing.T) http.Handler {
	t.Helper()
	cfg := traefik_xff_fix.CreateConfig()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
	handler, err := traefik_xff_fix.New(context.Background(), next, cfg, "xff-fix")
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func mustNewRequest(t *testing.T, xffValue string) *http.Request {
	t.Helper()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	if xffValue != "" {
		req.Header.Set("X-Forwarded-For", xffValue)
	}
	return req
}

func assertXFF(t *testing.T, req *http.Request, expected string) {
	t.Helper()
	got := req.Header.Get("X-Forwarded-For")
	if got != expected {
		t.Errorf("X-Forwarded-For: got %q, want %q", got, expected)
	}
}

func assertXFFDeleted(t *testing.T, req *http.Request) {
	t.Helper()
	values, ok := req.Header["X-Forwarded-For"]
	if ok {
		t.Errorf("X-Forwarded-For should be removed, got %q", values)
	}
}

func assertRemoteAddrHost(t *testing.T, req *http.Request, expectedHost string) {
	t.Helper()
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		t.Fatalf("RemoteAddr should be host:port, got %q: %v", req.RemoteAddr, err)
	}
	if host != expectedHost {
		t.Errorf("RemoteAddr host: got %q, want %q", host, expectedHost)
	}
}

package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var ErrUnknownService = errors.New("unknown service")

type Checker struct {
	checks      map[string]func(context.Context) error
	connections []*grpc.ClientConn
	http        *http.Client
}

// New creates reusable clients for configured targets. grpc:// targets use the
// standard health service; HTTP(S) targets must point to a health endpoint.
func New(targets map[string]string) (*Checker, error) {
	c := &Checker{
		checks: make(map[string]func(context.Context) error),
		http: &http.Client{
			Transport:     http.DefaultTransport.(*http.Transport).Clone(),
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
	}
	for name, target := range targets {
		u, err := url.Parse(target)
		if err != nil || name == "" || name == "backend" || u == nil || u.Host == "" || u.User != nil || u.Fragment != "" {
			c.Close()
			return nil, fmt.Errorf("invalid health target for %q", name)
		}
		switch u.Scheme {
		case "grpc":
			if u.Path != "" || u.RawQuery != "" {
				c.Close()
				return nil, fmt.Errorf("gRPC health target for %q must contain only host and port", name)
			}
			conn, err := grpc.NewClient(u.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				c.Close()
				return nil, fmt.Errorf("create health client for %q: %w", name, err)
			}
			c.connections = append(c.connections, conn)
			client := healthpb.NewHealthClient(conn)
			c.checks[name] = func(ctx context.Context) error {
				response, err := client.Check(ctx, &healthpb.HealthCheckRequest{})
				if err != nil {
					return err
				}
				if response.GetStatus() != healthpb.HealthCheckResponse_SERVING {
					return fmt.Errorf("gRPC health status: %s", response.GetStatus())
				}
				return nil
			}
		case "http", "https":
			c.checks[name] = func(ctx context.Context) error {
				request, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
				if err != nil {
					return err
				}
				response, err := c.http.Do(request)
				if err != nil {
					return err
				}
				defer response.Body.Close()
				// Drain small health responses to allow connection reuse.
				if _, err := io.Copy(io.Discard, io.LimitReader(response.Body, 4096)); err != nil {
					return err
				}
				if response.StatusCode != http.StatusOK {
					return fmt.Errorf("HTTP health status: %d", response.StatusCode)
				}
				return nil
			}
		default:
			c.Close()
			return nil, fmt.Errorf("unsupported health protocol for %q", name)
		}
	}
	return c, nil
}

func (c *Checker) Check(ctx context.Context, service string) error {
	check, ok := c.checks[service]
	if !ok {
		return ErrUnknownService
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return check(ctx)
}

func (c *Checker) Close() {
	for _, conn := range c.connections {
		conn.Close()
	}
	c.http.CloseIdleConnections()
}

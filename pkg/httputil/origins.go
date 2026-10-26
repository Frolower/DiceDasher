package httputil

import (
	"fmt"
	"net/url"
	"strings"
)

// ParseAllowedOrigins validates a comma-separated list of exact HTTP origins.
// An unset or whitespace-only value uses defaults. Empty list entries are invalid.
func ParseAllowedOrigins(raw, defaults string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		raw = defaults
	}
	origins := strings.Split(raw, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
		parsed, err := url.Parse(origins[i])
		if err != nil || parsed.Hostname() == "" ||
			(parsed.Scheme != "http" && parsed.Scheme != "https") ||
			parsed.User != nil || parsed.Path != "" ||
			parsed.RawQuery != "" || parsed.ForceQuery ||
			strings.Contains(origins[i], "#") || strings.Contains(parsed.Host, "*") {
			return nil, fmt.Errorf("ALLOWED_ORIGINS contains invalid origin %q: expected an http(s) origin without path", origins[i])
		}
	}
	return origins, nil
}

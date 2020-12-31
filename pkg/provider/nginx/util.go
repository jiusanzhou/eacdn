package nginx

import (
	"fmt"
	"strings"
)

func parseVersion(s string) (string, error) {
	parts := strings.Split(strings.TrimSpace(s), "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("nginx version layout not corrent: %s", s)
	}
	return parts[1], nil
}

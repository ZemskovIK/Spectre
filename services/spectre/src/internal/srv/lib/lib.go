package lib

import (
	"strconv"
	"strings"
)

const (
	UNK_NAME     = "" // ! TODO : in cfg
	ADMIN_ALEVEL = 6  // max level in system
)

// GetID extracts the numeric ID from a URL pattern after a given point.
// Returns the string ID, integer ID, and error if conversion fails.
func GetID(point, pattern string) (string, int, error) {
	var sid string
	if parts := strings.Split(pattern, point); len(parts) > 1 {
		sid = parts[1]
	}
	id, err := strconv.Atoi(sid)
	if err != nil {
		return sid, -1, err
	}

	return sid, id, nil
}

package lib

import (
	"strconv"
	"strings"
)

const (
	UNK_NAME     = ""
	ADMIN_ALEVEL = 6 // max level in system
)

// SplitName splits a full name string into first, middle, and last name parts.
func SplitName(name string) (string, string, string) {
	names := strings.Split(name, " ")
	switch len(names) {
	case 1:
		return names[0], UNK_NAME, UNK_NAME
	case 2:
		return names[0], names[1], UNK_NAME
	case 3:
		return names[0], names[1], names[2]
	default:
		return names[0], names[1], strings.Join(names[2:], " ")
	}
}

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

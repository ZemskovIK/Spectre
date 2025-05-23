package lib

import (
	"strconv"
	"strings"
)

const (
	UNK_NAME     = ""
	ADMIN_ALEVEL = 6
)

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

// GetID returns the last post-/ part of pattern slited by point
// ex: pattern '/api/point/1/2 point' 'GET /api/point/' : will return "1/2" -1 err
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

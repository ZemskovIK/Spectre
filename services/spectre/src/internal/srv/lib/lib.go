package lib

import (
	"spectre/internal/models"
	"spectre/pkg/logger"
	"strconv"
	"strings"
)

const (
	UNK_NAME     = "" // ! TODO : in cfg
	ADMIN_ALEVEL = 6  // max level in system

	GLOC_LIB = "/src/internal/srv/lib/lib.go/"
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

// ValidateLetter validates the letter and write logs in log.
// Returns error message if fails and ok status.
func ValidateLetter(
	letter models.Letter, log *logger.Logger,
) (string, bool) {
	loc := GLOC_LIB + "LetterIsValid()"

	if letter.Body == "" {
		log.Warnf("%s: body cannot be empty", loc)
		return "body cannot be empty!", false
	}
	if letter.Author == "" {
		log.Infof("%s: author is empty, setting to 'unknown'", loc)
		letter.Author = "unknown"
	}
	if letter.FoundIn == "" {
		log.Infof("%s: foundIn is empty, setting to 'unknown'", loc)
		letter.FoundIn = "unknown"
	}

	return "", true
}

package lib

import "strings"

const (
	UNK_NAME = ""
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

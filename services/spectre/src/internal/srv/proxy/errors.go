package proxy

import "fmt"

var (
	errBadStatusCode = func(sc int) error {
		return fmt.Errorf("bad status code from proxy: %d", sc)
	}
)

package daemon

import (
	"errors"
	"os"
)

// ErrNotRunning is returned when no daemon socket is reachable at the given path.
// cmd/* surfaces this as the "run aw on first" hint.
var ErrNotRunning = errors.New("aw daemon not running")

// Ping checks whether a daemon is reachable at sockPath.
// Returns ErrNotRunning when the socket is missing or refuses connections.
func Ping(sockPath string) error {
	if _, err := os.Stat(sockPath); err != nil {
		if os.IsNotExist(err) {
			return ErrNotRunning
		}
		return err
	}
	// TODO: actually dial unix sock + send OpPing, treat ECONNREFUSED as ErrNotRunning
	return nil
}

// Call sends a single Request to the daemon and returns its Response.
func Call(sockPath string, req Request) (Response, error) {
	if err := Ping(sockPath); err != nil {
		return Response{}, err
	}
	// TODO: dial unix sock, write JSON request, read JSON response
	return Response{OK: true}, nil
}

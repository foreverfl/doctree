package daemon

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"syscall"
	"time"
)

// ErrNotRunning is returned when no daemon socket is reachable at the given path.
// cmd/* surfaces this as the "run gitt on first" hint.
var ErrNotRunning = errors.New("gitt daemon not running")

const (
	dialTimeout = 2 * time.Second
	rwTimeout   = 5 * time.Second
)

// Ping checks whether a daemon is reachable at sockPath by issuing OpPing.
// Returns ErrNotRunning when the socket is missing or the dial is refused.
func Ping(sockPath string) error {
	if _, err := os.Stat(sockPath); err != nil {
		if os.IsNotExist(err) {
			return ErrNotRunning
		}
		return err
	}
	resp, err := Call(sockPath, Request{Op: OpPing})
	if err != nil {
		return err
	}
	if !resp.OK {
		return errors.New("daemon ping rejected: " + resp.Error)
	}
	return nil
}

// Call sends a single Request to the daemon and returns its Response.
// ECONNREFUSED / ENOENT on dial map to ErrNotRunning so callers can give the
// same "run gitt on first" hint regardless of which failure mode hit.
func Call(sockPath string, req Request) (Response, error) {
	conn, err := net.DialTimeout("unix", sockPath, dialTimeout)
	if err != nil {
		if isNotRunning(err) {
			return Response{}, ErrNotRunning
		}
		return Response{}, err
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(rwTimeout)); err != nil {
		return Response{}, err
	}
	if err := json.NewEncoder(conn).Encode(req); err != nil {
		return Response{}, err
	}
	var resp Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return Response{}, err
	}
	return resp, nil
}

func isNotRunning(err error) bool {
	if errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ENOENT) {
		return true
	}
	var sysErr *os.SyscallError
	if errors.As(err, &sysErr) {
		return errors.Is(sysErr.Err, syscall.ECONNREFUSED) || errors.Is(sysErr.Err, syscall.ENOENT)
	}
	return false
}

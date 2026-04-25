package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type handler func(req Request) Response

type server struct {
	listener     net.Listener
	handlers     map[Op]handler
	shutdown     chan struct{}
	shutdownOnce sync.Once
	wg           sync.WaitGroup
}

// Run starts the daemon: opens the SQLite store at dbPath (TODO commit 3) and
// listens on sockPath. Blocks until OpShutdown / SIGTERM / SIGINT.
func Run(sockPath, dbPath string) error {
	_ = dbPath // TODO commit 3: store.Open(dbPath)

	if err := os.Remove(sockPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("clear stale sock: %w", err)
	}

	ln, err := net.Listen("unix", sockPath)
	if err != nil {
		return fmt.Errorf("listen %s: %w", sockPath, err)
	}

	server := &server{
		listener: ln,
		shutdown: make(chan struct{}),
	}
	server.handlers = map[Op]handler{
		OpPing: func(_ Request) Response { return Response{OK: true} },
		OpShutdown: func(_ Request) Response {
			// Close listener first so accept loop exits, then signal Run to
			// finish. Done in a goroutine so this handler can still return its
			// response on the same connection before the conn is torn down.
			go server.shutdownOnce.Do(func() { close(server.shutdown) })
			return Response{OK: true}
		},
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	go server.acceptLoop()

	select {
	case <-server.shutdown:
	case <-sigCh:
	}

	_ = ln.Close()
	server.wg.Wait()
	_ = os.Remove(sockPath)
	return nil
}

func (server *server) acceptLoop() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Fprintln(os.Stderr, "doctree: accept:", err)
			continue
		}
		server.wg.Add(1)
		go func() {
			defer server.wg.Done()
			server.handleConn(conn)
		}()
	}
}

func (server *server) handleConn(conn net.Conn) {
	defer conn.Close()

	var req Request
	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		_ = json.NewEncoder(conn).Encode(Response{OK: false, Error: "decode: " + err.Error()})
		return
	}

	h, ok := server.handlers[req.Op]
	if !ok {
		_ = json.NewEncoder(conn).Encode(Response{OK: false, Error: "unknown op: " + string(req.Op)})
		return
	}
	_ = json.NewEncoder(conn).Encode(h(req))
}

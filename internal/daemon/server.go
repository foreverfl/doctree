package daemon

// Run starts the daemon: opens the SQLite store at dbPath and listens on sockPath.
// All requests are serialized through the store; this is the chokepoint that
// prevents port-allocation races across concurrent aw invocations.
func Run(sockPath, dbPath string) error {
	// TODO: open store.Store(dbPath)
	// TODO: net.Listen("unix", sockPath); defer remove
	// TODO: accept loop -> per-conn goroutine -> JSON request -> dispatch -> JSON response
	// TODO: graceful shutdown on OpShutdown / SIGTERM
	return nil
}

package observability

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/pprof"
)

// StartPprofServer starts an HTTP server serving pprof endpoints on the given address.
// This allows runtime profiling (CPU, Heap, Mutex, Block).
func StartPprofServer(ctx context.Context, addr string) error {
	mux := http.NewServeMux()

	// Register all pprof endpoints
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		slog.Info("observability: stopping pprof server")
		srv.Shutdown(context.Background())
	}()

	slog.Info("observability: starting pprof server", "addr", addr)
	return srv.ListenAndServe()
}

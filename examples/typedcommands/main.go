package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/containeroo/tinyflags"
)

// main demonstrates deferred command handlers with parsed flag values.
func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"admin", "users", "--name=bob", "--verbose"}
	}

	if err := run(context.Background(), args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// run parses flags, prepares cancellation, and executes the selected command.
func run(ctx context.Context, args []string) error {
	runner, err := parseFlags(args)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	return runner.Run(ctx)
}

// parseFlags registers commands and returns one deferred runner for the selected branch.
func parseFlags(args []string) (tinyflags.Runner, error) {
	app := tinyflags.NewCommand("app", tinyflags.ContinueOnError)
	verbose := app.Globals().Bool("verbose", false, "Enable verbose logging").Value()
	app.Run(runRoot, verbose)

	serve := app.Command("serve", "Run the HTTP server")
	listenAddr := serve.String("listen-addr", "127.0.0.1", "Listen address").Value()
	port := serve.Int("port", 8080, "HTTP port").Value()
	readHeaderTimeout := serve.Duration("read-header-timeout", 5*time.Second, "Read header timeout").Value()
	writeTimeout := serve.Duration("write-timeout", 10*time.Second, "Write timeout").Value()
	idleTimeout := serve.Duration("idle-timeout", 30*time.Second, "Idle timeout").Value()
	shutdownTimeout := serve.Duration("shutdown-timeout", 20*time.Second, "Graceful shutdown timeout").Value()
	serve.Run(runServer, listenAddr, verbose, port, readHeaderTimeout, writeTimeout, idleTimeout, shutdownTimeout)

	admin := app.Command("admin", "Administrative commands")
	users := admin.Command("users", "Manage user accounts")
	name := users.String("name", "alice", "User name").Value()
	users.Run(runAdminUsers, verbose, name)

	return app.ParseRunner(args)
}

// runRoot handles the root command when no subcommand is selected.
func runRoot(_ context.Context, verbose bool) error {
	fmt.Printf("root command selected (verbose=%t)\n", verbose)
	return nil
}

// runServer starts one real HTTP server and shuts it down when the context is canceled.
func runServer(
	ctx context.Context,
	listenAddr string,
	verbose bool,
	port int,
	readHeaderTimeout time.Duration,
	writeTimeout time.Duration,
	idleTimeout time.Duration,
	shutdownTimeout time.Duration,
) error {
	addr := net.JoinHostPort(listenAddr, strconv.Itoa(port))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "tinyflags server listening on %s\n", addr)
	})

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	if verbose {
		fmt.Printf("starting server at %s\n", addr)
		fmt.Printf("timeouts: read-header=%s write=%s idle=%s shutdown=%s\n",
			readHeaderTimeout,
			writeTimeout,
			idleTimeout,
			shutdownTimeout,
		)
	}

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// runAdminUsers demonstrates a nested command with parsed flag values.
func runAdminUsers(_ context.Context, verbose bool, name string) error {
	fmt.Printf("managing user %q\n", name)
	fmt.Printf("verbose: %t\n", verbose)
	return nil
}

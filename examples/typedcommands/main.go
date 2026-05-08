package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
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
	listenAddr := serve.String("listen-addr", "127.0.0.1:8080", "Listen address").Value()
	readHeaderTimeout := serve.Duration("read-header-timeout", 5*time.Second, "Read header timeout").Value()
	writeTimeout := serve.Duration("write-timeout", 10*time.Second, "Write timeout").Value()
	idleTimeout := serve.Duration("idle-timeout", 30*time.Second, "Idle timeout").Value()
	shutdownTimeout := serve.Duration("shutdown-timeout", 20*time.Second, "Graceful shutdown timeout").Value()

	serve.BuildCommand(func() tinyflags.Runner {
		return serverRunner{
			listenAddr:        *listenAddr,
			verbose:           *verbose,
			readHeaderTimeout: *readHeaderTimeout,
			writeTimeout:      *writeTimeout,
			idleTimeout:       *idleTimeout,
			shutdownTimeout:   *shutdownTimeout,
			args:              append([]string(nil), serve.Args()...),
		}
	})

	runCmd := app.Command("run", "Run a task and show positional arguments")
	runMode := runCmd.String("mode", "demo", "Execution mode").Value()
	runCmd.Run(
		func(_ context.Context, verbose bool, mode string) error {
			return runTask(verbose, mode, append([]string(nil), runCmd.Args()...))
		},
		verbose,
		runMode,
	)

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

type serverRunner struct {
	listenAddr        string
	verbose           bool
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
	shutdownTimeout   time.Duration
	args              []string
}

func (r serverRunner) Run(ctx context.Context) error {
	return runServer(
		ctx,
		r.listenAddr,
		r.verbose,
		r.readHeaderTimeout,
		r.writeTimeout,
		r.idleTimeout,
		r.shutdownTimeout,
		r.args,
	)
}

// runServer starts one real HTTP server and shuts it down when the context is canceled.
func runServer(
	ctx context.Context,
	listenAddr string,
	verbose bool,
	readHeaderTimeout time.Duration,
	writeTimeout time.Duration,
	idleTimeout time.Duration,
	shutdownTimeout time.Duration,
	args []string,
) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "tinyflags server listening on %s\n", listenAddr)
		fmt.Fprintf(w, "args: %v\n", args)
	})

	server := &http.Server{
		Addr:              listenAddr,
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
		fmt.Printf("starting server at %s\n", listenAddr)
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

func runTask(verbose bool, mode string, args []string) error {
	fmt.Printf("run mode: %s\n", mode)
	fmt.Printf("verbose: %t\n", verbose)
	fmt.Printf("args: %v\n", args)
	return nil
}

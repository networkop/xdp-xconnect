package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"os/signal"

	"github.com/networkop/xdp-xconnect/pkg/xdp"
)

var (
	configFlag = flag.String("conf", "", "xdp-xconnect configuration file (YAML)")
	version    = flag.Bool("v", false, "Display version")
)

func setupSigHandlers(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	go func() {
		sig := <-sigs
		log.Printf("Received syscall:%+v", sig)
		cancel()
	}()

}

func printVersion(gitCommit string) {
	fmt.Printf("Version: %s\n", gitCommit)
	os.Exit(0)
}

// Run the xdp-xconnect app
func Run(commit string) error {
	fmt.Println("eBPF x-connect")

	ctx, cancel := context.WithCancel(context.Background())

	setupSigHandlers(cancel)

	flag.Parse()

	if *version {
		printVersion(commit)
	}

	cfg, err := newFromFile(*configFlag)
	if err != nil {
		return fmt.Errorf("Parsing configuration file: %s", err)
	}

	updateCh := make(chan map[string]string, 1)
	go configWatcher(*configFlag, updateCh)

	app, err := xdp.NewXconnectApp(cfg.Links)
	if err != nil {
		return fmt.Errorf("Loading eBPF: %s", err)
	}

	app.Launch(ctx, updateCh)

	return nil
}

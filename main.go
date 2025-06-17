package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	var (
		name  string
		addr  string
		debug bool
	)
	flag.StringVar(&name, "name", "", "set name of the peer")
	flag.StringVar(&addr, "addr", "", "set the address of your system, example 192.168.0.103:3030")
	flag.BoolVar(&debug, "debug", false, "enable debug mode")
	flag.Parse()

	if name == "" {
		fmt.Println("-name flag is required to be set")
		os.Exit(1)
	}
	if addr == "" {
		fmt.Println("-addr flag is required to be set")
		os.Exit(1)
	}
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if debug {
		opts.Level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	discovery := newDiscovery(name, addr, logger)
	manager := newManager(addr, logger)

	fmt.Println("Welcome to your peering relationship!")
	if err := scanFromKeyboard(manager, discovery); err != nil {
		logger.With("err", err.Error()).Error("error scanning input from keyboard")
		os.Exit(1)
	}

}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func main() {
	configPath := flag.String("config", "", "path to config file (optional)")
	learnMode := flag.Bool("learn", false, "run in learn mode: record current ports as baseline and exit")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("portwatch: failed to load config: %v", err)
	}

	bl, err := baseline.Load(cfg.BaselinePath)
	if err != nil {
		log.Printf("portwatch: no existing baseline found, starting fresh: %v", err)
		bl = baseline.New()
	}

	sc, err := scanner.NewProcScanner()
	if err != nil {
		log.Fatalf("portwatch: failed to create scanner: %v", err)
	}

	if *learnMode {
		runLearnMode(sc, bl, cfg)
		return
	}

	alerter := alert.NewAlerter(os.Stdout)
	mon := monitor.New(sc, bl, alerter, cfg)

	fmt.Println("portwatch: starting daemon...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := mon.Run(ctx); err != nil {
		log.Fatalf("portwatch: monitor exited with error: %v", err)
	}

	fmt.Println("portwatch: shutting down.")
}

func runLearnMode(sc scanner.Scanner, bl *baseline.Baseline, cfg *config.Config) {
	states, err := sc.Scan()
	if err != nil {
		log.Fatalf("portwatch: scan failed: %v", err)
	}
	for _, s := range states {
		bl.Add(s)
	}
	if err := bl.Save(cfg.BaselinePath); err != nil {
		log.Fatalf("portwatch: failed to save baseline: %v", err)
	}
	fmt.Printf("portwatch: learned %d port(s), baseline saved to %s\n", len(states), cfg.BaselinePath)
}

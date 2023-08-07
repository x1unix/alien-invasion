package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/x1unix/alien-invasion/internal/invasion"
	"github.com/x1unix/alien-invasion/internal/mapfile"
	"go.uber.org/zap"
)

func main() {
	if err := mainErr(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func mainErr() error {
	var (
		aliensCount int
		moveLimit   uint
		verbose     bool
		file        string
	)

	flag.StringVar(&file, "f", "", "map file name")
	flag.IntVar(&aliensCount, "c", 0, "aliens count")
	flag.UintVar(&moveLimit, "l", 10000, "game moves limit")
	flag.BoolVar(&verbose, "v", false, "enable game progress logging")
	flag.Parse()

	if verbose {
		if err := initDebugLogger(); err != nil {
			return err
		}
	}

	if file == "" {
		return errors.New("missing file name")
	}

	if aliensCount == 0 {
		return errors.New("missing aliens count")
	}

	s, err := mapfile.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read map: %w", err)
	}

	for aliensCount > len(s.Cities) {
		return fmt.Errorf(
			"number of aliens (%d) cannot exceed number of cities (%d)",
			aliensCount, len(s.Cities),
		)
	}

	game := invasion.NewGame(aliensCount, moveLimit, s.Cities)
	for game.Tick() {
	}

	return nil
}

func initDebugLogger() error {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"stderr"}
	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	return nil
}

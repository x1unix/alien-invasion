package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/x1unix/alien-invasion/internal/invasion"
	"github.com/x1unix/alien-invasion/internal/mapfile"
)

func main() {
	if err := mainErr(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func mainErr() error {
	var (
		aliensCount int
		moveLimit   uint
		file        string
	)

	flag.StringVar(&file, "f", "", "map file name")
	flag.IntVar(&aliensCount, "c", 0, "aliens count")
	flag.UintVar(&moveLimit, "l", 10000, "game moves limit")
	flag.Parse()

	if file == "" {
		return errors.New("missing file name")
	}

	if aliensCount == 0 {
		return errors.New("missing aliens count")
	}

	s, err := mapfile.ReadFile(file)
	if err != nil {
		return err
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

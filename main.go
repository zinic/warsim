package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"
)

func writeNewState() {
	fmt.Printf("No state file exists, creating a new one.\n")

	newState := &Simulation{
		Step:     0,
		StateDir: "state",
	}

	if err := newState.Write(); err != nil {
		panic(fmt.Sprintf("Failed to Write new state file: %v.", err))
	}

	os.Exit(0)
}

func checkStateFile() {
	if _, err := os.Stat(path.Join("state", stateFilename)); err != nil {
		if !os.IsNotExist(err) {

			panic(fmt.Sprintf("Faild to open state file: %v.", err))
		}

		writeNewState()
	}
}

func launchStdinReader(stdinC chan string) {
	reader := bufio.NewReader(os.Stdin)

	go func() {
		for {
			if text, err := reader.ReadString('\n'); err != nil {
				os.Exit(0)
			} else {
				stdinC <- text
			}
		}
	}()
}

func main() {
	rand.Seed(time.Now().Unix())
	checkStateFile()

	//stdinC := make(chan string)
	//launchStdinReader(stdinC)
	if simulation, err := LoadSimulation("state"); err != nil {
		panic(fmt.Sprintf("Failed to load state: %v.", err))
	} else if err := simulation.Turn(); err != nil {
		fmt.Printf("Error executing simulation: %v.", err)
	}
}

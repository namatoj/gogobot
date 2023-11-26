package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	GenmoveBlack = "genmove black\n"
	GenmoveWhite = "genmove white\n"
)

func generateMove(blackToPlay bool, stdin io.WriteCloser) {
	command := GenmoveWhite
	if blackToPlay {
		command = GenmoveBlack
	}
	io.WriteString(stdin, command)

}

func endGameSession(stdin io.WriteCloser) {
	io.WriteString(stdin, "showboard\n")
	io.WriteString(stdin, "final_score\n")
	io.WriteString(stdin, "quit\n")
}

func gameSetup(stdin io.WriteCloser, scanner *bufio.Scanner) {
	io.WriteString(stdin, "komi 6.5\n")
	io.WriteString(stdin, "boardsize 9\n")
	scanner.Scan()
	scanner.Scan()
	scanner.Scan()
}

func playGo() {
	cmd := exec.Command("/usr/games/gnugo", "--mode", "gtp")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Add the above to the gameSetup function
	scanner := bufio.NewScanner(stdout)
	gameSetup(stdin, scanner)

	blackToPlay := true
	moves := 0
	consecutivePasses := 0

	for scanner.Scan() {
		text := scanner.Text()

		if strings.Contains(text, "B+") {
			fmt.Println("Black wins", text)
			break
		}
		if strings.Contains(text, "W+") {
			fmt.Println("White wins", text)
			break
		}

		matchIsOver := consecutivePasses >= 2
		if strings.HasPrefix(text, "= ") && moves > 0 && !matchIsOver {
			if blackToPlay {
				fmt.Println("Move", moves, ": White plays", text)
			} else {
				fmt.Println("Move", moves, ": Black plays", text)
			}
			if strings.Contains(text, "PASS") {
				consecutivePasses++

			} else {
				consecutivePasses = 0
			}
		} else {
			if len(text) > 0 {
				fmt.Println(text)
			}
		}

		isReadyForInput := len(text) == 0
		if isReadyForInput {
			if !matchIsOver {
				generateMove(blackToPlay, stdin)
				blackToPlay = !blackToPlay
				moves++
			} else {
				endGameSession(stdin)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func main() {
	playGo()
}

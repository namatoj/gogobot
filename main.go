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

	scanner := bufio.NewScanner(stdout)
	io.WriteString(stdin, "komi 6.5\n")
	io.WriteString(stdin, "boardsize 9\n")
	scanner.Scan()
	scanner.Scan()

	blackToPlay := true
	genmoveBlack := "genmove black\n"
	genmoveWhite := "genmove white\n"
	moves := 0
	command := genmoveBlack
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

		if strings.HasPrefix(text, "= ") && moves > 0 {
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

		if len(text) == 0 {
			if consecutivePasses < 2 {
				if blackToPlay {
					command = genmoveBlack
					blackToPlay = false
				} else {
					command = genmoveWhite
					blackToPlay = true
				}
				io.WriteString(stdin, command)
				moves++
			} else {
				io.WriteString(stdin, "showboard\n")
				io.WriteString(stdin, "final_score\n")
				io.WriteString(stdin, "quit\n")
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

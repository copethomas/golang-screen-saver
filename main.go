/*

A Really Basic Terminal Scren Saver Written in Go
Created by Thomas Cope

*/

package main

import (
	tm "github.com/buger/goterm"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Dot struct {
	eaten    bool
	isTarget bool
	x        int
	y        int
}

type Pac struct {
	score int
	nom   bool
	x     int
	y     int
}

func getRandom(min, max int) int {
	return rand.Intn(max-min) + min
}

func drawDots(dots []Dot) {
	tm.Clear()
	for i := range dots {
		if dots[i].eaten == false {
			tm.MoveCursor(dots[i].x, dots[i].y)
			if dots[i].isTarget {
				tm.Print("#")
			} else {
				tm.Print("*")
			}
		}
	}
}

func findTargetDot(dots []Dot, pacman Pac) (int, int, int) {
	dis := 99999.9
	target := 1
	tdis := 99999.9
	for i := range dots {
		if dots[i].eaten == false {
			tdis = math.Sqrt(math.Pow(float64(dots[i].x-pacman.x), 2) + math.Pow(float64(dots[i].y-pacman.y), 2))
			if tdis < dis {
				dis = tdis
				target = i
			}
		}
	}
	return dots[target].y, dots[target].x, target
}

func main() {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-exitSignal
		tm.Clear()
		tm.Flush()
		os.Exit(0)
	}()
	rand.Seed(time.Now().Unix()) //Don't really need a cryptographicly secure random numbers for a screen saver.
	dots := []Dot{}
	for i := 0; i < getRandom(50, 100); i++ {
		tmp := Dot{eaten: false, x: (getRandom(1, tm.Width())), y: (getRandom(1, tm.Height()))}
		dots = append(dots, tmp)
	}
	pacman := Pac{score: 0, x: (getRandom(1, tm.Width())), y: (getRandom(1, tm.Height())), nom: false}
	ty, tx, t := findTargetDot(dots, pacman)
	dots[t].isTarget = true
	for {
		drawDots(dots)
		tm.MoveCursor(pacman.x, pacman.y)
		if pacman.nom {
			tm.Print("@")
			pacman.nom = false
		} else {
			tm.Print("0")
			pacman.nom = true
		}
		if ty == pacman.y && tx == pacman.x {
			dots[t].eaten = true
			ty, tx, t = findTargetDot(dots, pacman)
			dots[t].isTarget = true
			pacman.score++
		}
		if tx > pacman.x {
			pacman.x++
		} else if tx < pacman.x {
			pacman.x--
		}
		if ty > pacman.y {
			pacman.y++
		} else if ty < pacman.y {
			pacman.y--
		}
		if pacman.score == len(dots) {
			pacman.score = 0
			dots = []Dot{}
			for i := 0; i < getRandom(50, 100); i++ {
				tmp := Dot{eaten: false, x: (getRandom(1, tm.Width())), y: (getRandom(1, tm.Height()))}
				dots = append(dots, tmp)
			}
			ty, tx, t = findTargetDot(dots, pacman)
		}
		tm.MoveCursor(1, 1)
		tm.Flush()
		time.Sleep(time.Second / 5)
	}
}

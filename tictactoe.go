// The MIT License (MIT)
//
// Copyright (c) 2017 Fredy Wijaya
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
)

type coordinate struct {
	x int
	y int
}

const (
	startX = 1
	endX   = 13
	startY = 0
	endY   = 6
)

var (
	position = map[int]coordinate{
		1: {3, 1},
		2: {7, 1},
		3: {11, 1},
		4: {3, 3},
		5: {7, 3},
		6: {11, 3},
		7: {3, 5},
		8: {7, 5},
		9: {11, 5},
	}
)

func drawText(x, y int, text string) {
	colorDefault := termbox.ColorDefault
	for _, ch := range text {
		termbox.SetCell(x, y, ch, colorDefault, colorDefault)
		x++
	}
}

func drawBoard() {
	colorDefault := termbox.ColorDefault

	// 4 horizontal lines
	for line := startY; line <= endY; line += 2 {
		for x := startX; x <= endX; x++ {
			c := '\u2550'
			termbox.SetCell(x, line, c, colorDefault, colorDefault)
		}
	}

	// 4 vertical lines
	for line := startX; line <= endX; line += 4 {
		for y := startY; y <= endY; y++ {
			c := '\u2551'
			termbox.SetCell(line, y, c, colorDefault, colorDefault)
		}
	}

	// prettify the board
	for y := startY; y <= endY; y += 2 {
		for x := startX; x <= endX; x += 4 {
			var c rune
			if y == startY && x == startX {
				c = '\u2554'
			} else if y == startY && (x == startX+4 || x == startX+8) {
				c = '\u2566'
			} else if y == startY && x == endX {
				c = '\u2557'
			} else if (y == startY+2 || y == startY+4) && x == startX {
				c = '\u2560'
			} else if (y == startY+2 || y == startY+4) && x == endX {
				c = '\u2563'
			} else if y == endY && x == startX {
				c = '\u255A'
			} else if y == endY && (x == startX+4 || x == startX+8) {
				c = '\u2569'
			} else if y == endY && x == endX {
				c = '\u255D'
			} else {
				c = '\u256C'
			}
			termbox.SetCell(x, y, c, colorDefault, colorDefault)
		}
	}

	termbox.SetCell(position[1].x, position[1].y, 'O', colorDefault, colorDefault)
	termbox.SetCell(position[2].x, position[2].y, 'X', colorDefault, colorDefault)
	termbox.SetCell(position[3].x, position[3].y, 'X', colorDefault, colorDefault)
	termbox.SetCell(position[4].x, position[4].y, 'O', colorDefault, colorDefault)
	termbox.SetCell(position[5].x, position[5].y, 'X', colorDefault, colorDefault)
	termbox.SetCell(position[6].x, position[6].y, 'O', colorDefault, colorDefault)
	termbox.SetCell(position[7].x, position[7].y, 'X', colorDefault, colorDefault)
	termbox.SetCell(position[8].x, position[8].y, 'X', colorDefault, colorDefault)
	termbox.SetCell(position[9].x, position[9].y, 'O', colorDefault, colorDefault)
}

func redrawAll() {
	colorDefault := termbox.ColorDefault
	termbox.Clear(colorDefault, colorDefault)

	drawBoard()
	drawText(15, 1, "Player 1's turn:")
	drawText(15, 2, "(1, 2, 3, 4, 5, 6, 7, 8, 9)")
	drawText(2, 7, "Created by Fredy Wijaya")

	termbox.Flush()
}

func runGame() {
	err := termbox.Init()
	if err != nil {
		errorAndExit(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	redrawAll()
exitGame:
	for {
		select {
		case ev := <-eventQueue:
			switch ev.Key {
			case termbox.KeyEsc:
				break exitGame
			}
		}
		redrawAll()
	}
}

func errorAndExit(message interface{}) {
	fmt.Println(message)
	os.Exit(1)
}

func main() {
	runGame()
}

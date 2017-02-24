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

func drawBoard() {
	colorDefault := termbox.ColorDefault

	// -------------
	// | O | X | O |
	// -------------
	// | O | X | O |
	// -------------
	// | O | X | O |
	// -------------

	// TODO: refactor this code

	// top
	for x := 1; x <= 13; x++ {
		var c rune
		if x == 1 {
			c = '\u2554'
		} else if x == 13 {
			c = '\u2557'
		} else {
			c = '\u2550'
		}
		termbox.SetCell(x, 1, c, colorDefault, colorDefault)
	}

	// left
	for y := 2; y <= 6; y++ {
		var c rune
		c = '\u2551'
		termbox.SetCell(1, y, c, colorDefault, colorDefault)
	}

	// bottom
	for x := 1; x <= 13; x++ {
		var c rune
		if x == 1 {
			c = '\u255A'
		} else if x == 13 {
			c = '\u255D'
		} else {
			c = '\u2550'
		}
		termbox.SetCell(x, 7, c, colorDefault, colorDefault)
	}

	// right
	for y := 2; y <= 6; y++ {
		var c rune
		c = '\u2551'
		termbox.SetCell(13, y, c, colorDefault, colorDefault)
	}

	// first horizontal line
	for x := 1; x <= 13; x++ {
		var c rune
		if x == 1 {
			c = '\u2560'
		} else if x == 13 {
			c = '\u2563'
		} else {
			c = '\u2550'
		}
		termbox.SetCell(x, 3, c, colorDefault, colorDefault)
	}

	// second horizontal line
	for x := 1; x <= 13; x++ {
		var c rune
		if x == 1 {
			c = '\u2560'
		} else if x == 13 {
			c = '\u2563'
		} else {
			c = '\u2550'
		}
		termbox.SetCell(x, 5, c, colorDefault, colorDefault)
	}

	// first vertical line
	for y := 1; y <= 7; y++ {
		var c rune
		if y == 1 {
			c = '\u2566'
		} else if y == 3 || y == 5 {
			c = '\u256C'
		} else if y == 7 {
			c = '\u2569'
		} else {
			c = '\u2551'
		}
		termbox.SetCell(5, y, c, colorDefault, colorDefault)
	}

	// second vertical line
	for y := 1; y <= 7; y++ {
		var c rune
		if y == 1 {
			c = '\u2566'
		} else if y == 3 || y == 5 {
			c = '\u256C'
		} else if y == 7 {
			c = '\u2569'
		} else {
			c = '\u2551'
		}
		termbox.SetCell(9, y, c, colorDefault, colorDefault)
	}

	termbox.SetCell(3, 2, 'O', colorDefault, colorDefault)
	termbox.SetCell(7, 2, 'X', colorDefault, colorDefault)
	termbox.SetCell(11, 2, 'X', colorDefault, colorDefault)
	termbox.SetCell(3, 4, 'O', colorDefault, colorDefault)
	termbox.SetCell(7, 4, 'X', colorDefault, colorDefault)
	termbox.SetCell(11, 4, 'O', colorDefault, colorDefault)
	termbox.SetCell(3, 6, 'X', colorDefault, colorDefault)
	termbox.SetCell(7, 6, 'X', colorDefault, colorDefault)
	termbox.SetCell(11, 6, 'O', colorDefault, colorDefault)
}

func redrawAll() {
	colorDefault := termbox.ColorDefault
	termbox.Clear(colorDefault, colorDefault)

	drawBoard()

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

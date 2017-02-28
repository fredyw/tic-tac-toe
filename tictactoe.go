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
	"net/http"
	"os"

	"github.com/nsf/termbox-go"
	"github.com/urfave/cli"
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

// Game is a struct to store game information.
type Game struct {
	Board [][]rune
}

var (
	position = map[string]coordinate{
		"0|0": {3, 1},
		"0|1": {7, 1},
		"0|2": {11, 1},
		"1|0": {3, 3},
		"1|1": {7, 3},
		"1|2": {11, 3},
		"2|0": {3, 5},
		"2|1": {7, 5},
		"2|2": {11, 5},
	}
	game = Game{
		// 3x3 board
		Board: [][]rune{
			{' ', ' ', ' '},
			{' ', ' ', ' '},
			{' ', ' ', ' '},
		},
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

	for i := range game.Board {
		for j := range game.Board[i] {
			symbol := game.Board[i][j]
			drawSymbol(fmt.Sprintf("%d|%d", i, j), symbol)
		}
	}
}

func drawSymbol(key string, symbol rune) {
	colorDefault := termbox.ColorDefault
	termbox.SetCell(position[key].x, position[key].y, symbol, colorDefault, colorDefault)
}

func setX(pos int) {
	row, col, err := getRowCol(pos)
	if err != nil {
		return
	}
	game.Board[row][col] = 'X'
}

func setO(pos int) {
	row, col, err := getRowCol(pos)
	if err != nil {
		return
	}
	game.Board[row][col] = 'O'
}

func getRowCol(pos int) (row, col int, err error) {
	if pos <= 0 || pos >= 10 {
		return 0, 0, fmt.Errorf("Invalid position: %d", pos)
	}
	row = (pos - 1) / 3
	col = (pos - 1) % 3
	return row, col, nil
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

func startGame() {
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
			if ev.Key == termbox.KeyEsc {
				break exitGame
			} else if ev.Ch == '1' {
				setX(1)
			} else if ev.Ch == '2' {
				setX(2)
			} else if ev.Ch == '3' {
				setX(3)
			} else if ev.Ch == '4' {
				setX(4)
			} else if ev.Ch == '5' {
				setX(5)
			} else if ev.Ch == '6' {
				setX(6)
			} else if ev.Ch == '7' {
				setX(7)
			} else if ev.Ch == '8' {
				setX(8)
			} else if ev.Ch == '9' {
				setX(9)
			}
		}
		redrawAll()
	}
}

func startServer(name string, port uint) error {
	http.HandleFunc("/tictactoe", func(w http.ResponseWriter, req *http.Request) {
		startGame()
	})
	fmt.Println("Waiting for Tic-Tac-Toe client to connect...")
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func startClient(name, host string, port uint) error {
	return nil
}

func errorAndExit(message interface{}) {
	fmt.Println(message)
	os.Exit(1)
}

func main() {
	app := cli.NewApp()
	app.Name = "tic-tac-toe"
	app.Commands = []cli.Command{
		{
			Name:  "client",
			Usage: "Tic-Tac-Toe client",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "Player name",
					Value: "Player 1",
				},
				cli.StringFlag{
					Name:  "host",
					Usage: "Host name",
					Value: "localhost",
				},
				cli.UintFlag{
					Name:  "port",
					Usage: "Port number",
					Value: 8888,
				},
			},
			Action: func(c *cli.Context) error {
				return startClient(c.String("name"), c.String("host"), c.Uint("port"))
			},
		},
		{
			Name:  "server",
			Usage: "Tic-Tac-Toe server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "Player name",
					Value: "Player 2",
				},
				cli.UintFlag{
					Name:  "port",
					Usage: "Port number",
					Value: 8888,
				},
			},
			Action: func(c *cli.Context) error {
				return startServer(c.String("name"), c.Uint("port"))
			},
		},
	}
	app.Run(os.Args)
}

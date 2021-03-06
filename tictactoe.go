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
	"os"
	"strings"

	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/nsf/termbox-go"
	"github.com/urfave/cli"
)

type coordinate struct {
	x int
	y int
}

const (
	startX   = 1
	endX     = 13
	startY   = 0
	endY     = 6
	apiPath  = "/tictactoe"
	yourTurn = "Your turn"
	nextTurn = "Another player's turn"
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
)

func drawText(x, y int, text string) {
	colorDefault := termbox.ColorDefault
	for _, ch := range text {
		termbox.SetCell(x, y, ch, colorDefault, colorDefault)
		x++
	}
}

func drawBoard(game *Game) {
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

func setSymbol(game *Game, pos int, symbol rune, updateFunc func()) {
	row, col, err := getRowCol(pos)
	if err != nil {
		return
	}
	if game.Board[row][col] != ' ' {
		return
	}

	game.Board[row][col] = symbol

	updateFunc()
}

func getRowCol(pos int) (int, int, error) {
	if pos <= 0 || pos >= 10 {
		return 0, 0, fmt.Errorf("Invalid position: %d", pos)
	}
	row := (pos - 1) / 3
	col := (pos - 1) % 3
	return row, col, nil
}

func drawAll(game *Game, status []string) {
	colorDefault := termbox.ColorDefault
	termbox.Clear(colorDefault, colorDefault)

	drawBoard(game)
	for i, msg := range status {
		drawText(15, startY+i+1, msg)
	}
	drawText(15, 5, "Press ESC to quit")
	drawText(2, 7, "Created by Fredy Wijaya")

	termbox.Flush()
}

func startGame(player uint, conn *websocket.Conn) {
	err := termbox.Init()
	if err != nil {
		errorAndExit(err)
	}
	defer termbox.Close()

	var symbol rune
	if player == 1 {
		symbol = 'X'
	} else {
		symbol = 'O'
	}

	game := &Game{
		// 3x3 board
		Board: [][]rune{
			{' ', ' ', ' '},
			{' ', ' ', ' '},
			{' ', ' ', ' '},
		},
	}
	drawAll(game, []string{})

	eventQueue := make(chan termbox.Event, 64)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	if player == 2 {
		drawAll(game, []string{nextTurn})
		conn.ReadJSON(&game)
		drawAll(game, []string{yourTurn, availablePositions(game)})
	} else {
		drawAll(game, []string{yourTurn, availablePositions(game)})
	}

	winner := ' '
	updateFunc := func() {
		conn.WriteJSON(game)
		drawAll(game, []string{nextTurn})
		end := endGame(game)
		if end != ' ' {
			winner = end
			return
		}
		conn.ReadJSON(&game)
		// make sure to drain every key event triggered while blocking to read JSON data
		for len(eventQueue) > 0 {
			<-eventQueue
		}
		drawAll(game, []string{yourTurn, availablePositions(game)})
		end = endGame(game)
		if end != ' ' {
			winner = end
			return
		}
	}

exitGame:
	for {
		select {
		case ev := <-eventQueue:
			if ev.Key == termbox.KeyEsc {
				break exitGame
			} else if ev.Ch == '1' {
				setSymbol(game, 1, symbol, updateFunc)
			} else if ev.Ch == '2' {
				setSymbol(game, 2, symbol, updateFunc)
			} else if ev.Ch == '3' {
				setSymbol(game, 3, symbol, updateFunc)
			} else if ev.Ch == '4' {
				setSymbol(game, 4, symbol, updateFunc)
			} else if ev.Ch == '5' {
				setSymbol(game, 5, symbol, updateFunc)
			} else if ev.Ch == '6' {
				setSymbol(game, 6, symbol, updateFunc)
			} else if ev.Ch == '7' {
				setSymbol(game, 7, symbol, updateFunc)
			} else if ev.Ch == '8' {
				setSymbol(game, 8, symbol, updateFunc)
			} else if ev.Ch == '9' {
				setSymbol(game, 9, symbol, updateFunc)
			}
			if winner != ' ' {
				break exitGame
			}
		}
	}

	if hasWinner(winner) {
		msg := ""
		if winner == 'D' {
			msg = "DRAW"
		} else if winner == symbol {
			msg = "YOU WIN"
		} else {
			msg = "YOU LOSE"
		}
		drawAll(game, []string{msg})
	endGame:
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					break endGame
				}
			}
		}
	}
}

func hasWinner(winner rune) bool {
	return winner == 'D' || winner == 'O' || winner == 'X'
}

func availablePositions(game *Game) string {
	numbers := []string{}
	number := 1
	for i := range game.Board {
		for j := range game.Board[i] {
			if game.Board[i][j] == ' ' {
				numbers = append(numbers, fmt.Sprintf("%d", number))
			}
			number++
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(numbers, ", "))
}

func endGame(game *Game) rune {
	blankCount := 0
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			current := game.Board[row][col]
			if current == ' ' {
				blankCount++
			}

			// top
			if row-2 >= 0 {
				found := true
				for r := row; r >= 0; r-- {
					if current != game.Board[r][col] {
						found = false
					}
				}
				if found {
					return current
				}
			}

			// top right
			if row-2 >= 0 && col+2 < 3 {
				found := true
				for r, c := row, col; r >= 0 && c < 3; r, c = r-1, c+1 {
					if current != game.Board[r][c] {
						found = false
					}
				}
				if found {
					return current
				}
			}

			// right
			if col+2 < 3 {
				found := true
				for c := col; c < 3; c++ {
					if current != game.Board[row][c] {
						found = false
					}
				}
				if found {
					return current
				}
			}

			// bottom right
			if row+2 < 3 && col+2 < 3 {
				found := true
				for r, c := row, col; r < 3 && c < 3; r, c = r+1, c+1 {
					if current != game.Board[r][c] {
						found = false
					}
				}
				if found {
					return current
				}
			}

			// bottom
			if row+2 < 3 {
				found := true
				for r := row; r < 3; r++ {
					if current != game.Board[r][col] {
						found = false
					}
				}
				if found {
					return current
				}
			}

			// bottom left
			if row+2 < 3 && col-2 >= 0 {
				found := true
				for r, c := row, col; r < 3 && c >= 0; r, c = r+1, c-1 {
					if current != game.Board[r][c] {
						found = false
					}
				}
				if found {
					return current
				}
			}

			// left
			if col-2 >= 0 {
				found := true
				for c := col; c >= 0; c-- {
					if current != game.Board[row][c] {
						found = false
					}
				}
				if found {
					return current
				}
			}

			// top left
			if row-2 >= 0 && col-2 >= 0 {
				found := true
				for r, c := row, col; r >= 0 && c >= 0; r, c = r-1, c-1 {
					if current != game.Board[r][c] {
						found = false
					}
				}
				if found {
					return current
				}
			}
		}
	}
	if blankCount == 0 {
		return 'D'
	}
	return ' '
}

func startServer(port uint) error {
	upgrader := websocket.Upgrader{}
	http.HandleFunc(apiPath, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		startGame(2, conn)
	})
	fmt.Println("Waiting for Tic-Tac-Toe client to connect...")
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func startClient(host string, port uint) error {
	fmt.Println("Connecting to Tic-Tac-Teo server...")
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", host, port), Path: apiPath}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	startGame(1, conn)
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
				return startClient(c.String("host"), c.Uint("port"))
			},
		},
		{
			Name:  "server",
			Usage: "Tic-Tac-Toe server",
			Flags: []cli.Flag{
				cli.UintFlag{
					Name:  "port",
					Usage: "Port number",
					Value: 8888,
				},
			},
			Action: func(c *cli.Context) error {
				return startServer(c.Uint("port"))
			},
		},
	}
	app.Run(os.Args)
}

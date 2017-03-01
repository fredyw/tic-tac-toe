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
	"runtime/debug"
	"testing"
)

func TestGetRowCol(t *testing.T) {
	row, col, err := getRowCol(1)
	assertEquals(t, 0, row)
	assertEquals(t, 0, col)
	assertNil(t, err)

	row, col, err = getRowCol(2)
	assertEquals(t, 0, row)
	assertEquals(t, 1, col)
	assertNil(t, err)

	row, col, err = getRowCol(3)
	assertEquals(t, 0, row)
	assertEquals(t, 2, col)
	assertNil(t, err)

	row, col, err = getRowCol(4)
	assertEquals(t, 1, row)
	assertEquals(t, 0, col)
	assertNil(t, err)

	row, col, err = getRowCol(5)
	assertEquals(t, 1, row)
	assertEquals(t, 1, col)
	assertNil(t, err)

	row, col, err = getRowCol(6)
	assertEquals(t, 1, row)
	assertEquals(t, 2, col)
	assertNil(t, err)

	row, col, err = getRowCol(7)
	assertEquals(t, 2, row)
	assertEquals(t, 0, col)
	assertNil(t, err)

	row, col, err = getRowCol(8)
	assertEquals(t, 2, row)
	assertEquals(t, 1, col)
	assertNil(t, err)

	row, col, err = getRowCol(9)
	assertEquals(t, 2, row)
	assertEquals(t, 2, col)
	assertNil(t, err)

	_, _, err = getRowCol(0)
	assertNotNil(t, err)

	_, _, err = getRowCol(10)
	assertNotNil(t, err)
}

func TestEndGame(t *testing.T) {
	game := &Game{
		Board: [][]rune{
			{' ', ' ', ' '},
			{' ', ' ', ' '},
			{' ', ' ', ' '},
		},
	}
	assertEquals(t, ' ', endGame(game))

	game = &Game{
		Board: [][]rune{
			{'X', ' ', ' '},
			{'X', 'O', 'X'},
			{'X', 'O', ' '},
		},
	}
	assertEquals(t, 'X', endGame(game))

	game = &Game{
		Board: [][]rune{
			{'X', ' ', ' '},
			{'O', 'X', ' '},
			{'X', ' ', 'X'},
		},
	}
	assertEquals(t, 'X', endGame(game))

	game = &Game{
		Board: [][]rune{
			{'O', 'O', 'O'},
			{'O', 'X', 'O'},
			{'X', 'O', 'X'},
		},
	}
	assertEquals(t, 'O', endGame(game))

	game = &Game{
		Board: [][]rune{
			{'O', 'X', 'O'},
			{'X', 'O', 'O'},
			{'X', 'X', 'X'},
		},
	}
	assertEquals(t, 'X', endGame(game))

	game = &Game{
		Board: [][]rune{
			{'O', 'X', 'O'},
			{'X', 'O', 'O'},
			{'O', 'X', 'X'},
		},
	}
	assertEquals(t, 'O', endGame(game))

	game = &Game{
		Board: [][]rune{
			{'O', 'X', 'O'},
			{'O', 'O', 'O'},
			{'X', 'O', 'X'},
		},
	}
	assertEquals(t, 'O', endGame(game))

	game = &Game{
		Board: [][]rune{
			{'X', 'X', 'X'},
			{'O', ' ', 'O'},
			{'X', 'O', 'X'},
		},
	}
	assertEquals(t, 'X', endGame(game))

	game = &Game{
		Board: [][]rune{
			{'O', 'X', 'O'},
			{'O', 'X', 'O'},
			{'X', 'O', 'X'},
		},
	}
	assertEquals(t, 'D', endGame(game))

	game = &Game{
		Board: [][]rune{
			{'O', ' ', 'X'},
			{' ', 'X', ' '},
			{'O', ' ', ' '},
		},
	}
	assertEquals(t, ' ', endGame(game))
}

func assertEquals(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		debug.PrintStack()
		t.Fail()
	}
}

func assertNil(t *testing.T, object interface{}) {
	if object != nil {
		debug.PrintStack()
		t.Fail()
	}
}

func assertNotNil(t *testing.T, object interface{}) {
	if object == nil {
		debug.PrintStack()
		t.Fail()
	}
}

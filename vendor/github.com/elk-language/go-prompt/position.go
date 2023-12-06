package prompt

import (
	istrings "github.com/elk-language/go-prompt/strings"
	"github.com/rivo/uniseg"
)

// Position stores the coordinates
// of a p.
//
// (0, 0) represents the top-left corner of the prompt,
// while (n, n) the bottom-right corner.
type Position struct {
	X istrings.Width
	Y int
}

// Join two positions and return a new position.
func (p Position) Join(other Position) Position {
	if other.Y == 0 {
		p.X += other.X
	} else {
		p.X = other.X
		p.Y += other.Y
	}
	return p
}

// Add two positions and return a new position.
func (p Position) Add(other Position) Position {
	return Position{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

// Subtract two positions and return a new position.
func (p Position) Subtract(other Position) Position {
	return Position{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
}

// positionAtEndOfString calculates the position of the
// p at the end of the given string.
func positionAtEndOfStringLine(str string, columns istrings.Width, line int) Position {
	var down int
	var right istrings.Width
	g := uniseg.NewGraphemes(str)

charLoop:
	for g.Next() {
		runes := g.Runes()

		if len(runes) == 1 && runes[0] == '\n' {
			if down == line {
				break charLoop
			}
			down++
			right = 0
		}

		right += istrings.Width(g.Width())
		if right > columns {
			if down == line {
				right = columns - 1
				break charLoop
			}
			right = istrings.Width(g.Width())
			down++
		}

	}

	return Position{
		X: right,
		Y: down,
	}
}

// positionAtEndOfString calculates the position
// at the end of the given string.
func positionAtEndOfString(str string, columns istrings.Width) Position {
	var down int
	var right istrings.Width
	g := uniseg.NewGraphemes(str)

	for g.Next() {
		runes := g.Runes()

		if len(runes) == 1 && runes[0] == '\n' {
			down++
			right = 0
			continue
		}

		right += istrings.Width(g.Width())
		if right == columns {
			right = 0
			down++
		}
	}

	return Position{
		X: right,
		Y: down,
	}
}

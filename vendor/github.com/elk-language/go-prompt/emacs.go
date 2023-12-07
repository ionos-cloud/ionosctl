package prompt

import (
	"github.com/elk-language/go-prompt/debug"
	istrings "github.com/elk-language/go-prompt/strings"
)

/*

========
PROGRESS
========

Moving the cursor
-----------------

* [x] Ctrl + a   Go to the beginning of the line (Home)
* [x] Ctrl + e   Go to the End of the line (End)
* [x] Ctrl + p   Previous command (Up arrow)
* [x] Ctrl + n   Next command (Down arrow)
* [x] Ctrl + f   Forward one character
* [x] Ctrl + b   Backward one character
* [x] Ctrl + xx  Toggle between the start of line and current cursor position

Editing
-------

* [x] Ctrl + L   Clear the Screen, similar to the clear command
* [x] Ctrl + d   Delete character under the cursor
* [x] Ctrl + h   Delete character before the cursor (Backspace)

* [x] Ctrl + w   Cut the Word before the cursor to the clipboard.
* [x] Ctrl + k   Cut the Line after the cursor to the clipboard.
* [x] Ctrl + u   Cut/delete the Line before the cursor to the clipboard.

* [ ] Ctrl + t   Swap the last two characters before the cursor (typo).
* [ ] Esc  + t   Swap the last two words before the cursor.

* [ ] ctrl + y   Paste the last thing to be cut (yank)
* [ ] ctrl + _   Undo

*/

var emacsKeyBindings = []KeyBind{
	// Go to the End of the line
	{
		Key: ControlE,
		Fn: func(p *Prompt) bool {
			return p.CursorRightRunes(
				istrings.RuneCountInString(p.buffer.Document().CurrentLineAfterCursor()),
			)
		},
	},
	// Go to the beginning of the line
	{
		Key: ControlA,
		Fn: func(p *Prompt) bool {
			return p.CursorLeftRunes(
				p.buffer.Document().FindStartOfFirstWordOfLine(),
			)
		},
	},
	// Cut the Line after the cursor
	{
		Key: ControlK,
		Fn: func(p *Prompt) bool {
			p.buffer.DeleteRunes(
				istrings.RuneCountInString(p.buffer.Document().CurrentLineAfterCursor()),
				p.renderer.col,
				p.renderer.row,
			)
			return true
		},
	},
	// Cut/delete the Line before the cursor
	{
		Key: ControlU,
		Fn: func(p *Prompt) bool {
			p.buffer.DeleteBeforeCursorRunes(
				istrings.RuneCountInString(p.buffer.Document().CurrentLineBeforeCursor()),
				p.renderer.col,
				p.renderer.row,
			)
			return true
		},
	},
	// Delete character under the cursor
	{
		Key: ControlD,
		Fn: func(p *Prompt) bool {
			if p.buffer.Text() != "" {
				p.buffer.Delete(1, p.renderer.col, p.renderer.row)
				return true
			}
			return false
		},
	},
	// Backspace
	{
		Key: ControlH,
		Fn: func(p *Prompt) bool {
			p.buffer.DeleteBeforeCursor(1, p.renderer.col, p.renderer.row)
			return true
		},
	},
	// Right allow: Forward one character
	{
		Key: ControlF,
		Fn: func(p *Prompt) bool {
			return p.CursorRight(1)
		},
	},
	// Alt Right allow: Forward one word
	{
		Key: AltRight,
		Fn: func(p *Prompt) bool {
			return p.CursorRightRunes(
				p.buffer.Document().FindRuneNumberUntilEndOfCurrentWord(),
			)
		},
	},
	// Left allow: Backward one character
	{
		Key: ControlB,
		Fn: func(p *Prompt) bool {
			return p.CursorLeft(1)
		},
	},
	// Alt Left allow: Backward one word
	{
		Key: AltLeft,
		Fn: func(p *Prompt) bool {
			return p.CursorLeftRunes(
				p.buffer.Document().FindRuneNumberUntilStartOfPreviousWord(),
			)
		},
	},
	// Cut the Word before the cursor.
	{
		Key: ControlW,
		Fn:  DeleteWordBeforeCursor,
	},
	{
		Key: AltBackspace,
		Fn:  DeleteWordBeforeCursor,
	},
	// Clear the Screen, similar to the clear command
	{
		Key: ControlL,
		Fn: func(p *Prompt) bool {
			consoleWriter.EraseScreen()
			consoleWriter.CursorGoTo(0, 0)
			debug.AssertNoError(consoleWriter.Flush())
			return true
		},
	},
}

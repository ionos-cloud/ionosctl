package prompt

import (
	istrings "github.com/elk-language/go-prompt/strings"
)

// GoLineEnd Go to the End of the line
func GoLineEnd(p *Prompt) bool {
	x := []rune(p.buffer.Document().TextAfterCursor())
	return p.CursorRightRunes(istrings.RuneNumber(len(x)))
}

// GoLineBeginning Go to the beginning of the line
func GoLineBeginning(p *Prompt) bool {
	x := []rune(p.buffer.Document().TextBeforeCursor())
	return p.CursorLeftRunes(istrings.RuneNumber(len(x)))
}

// DeleteChar Delete character under the cursor
func DeleteChar(p *Prompt) bool {
	p.buffer.Delete(1, p.renderer.col, p.renderer.row)
	return true
}

// DeleteBeforeChar Go to Backspace
func DeleteBeforeChar(p *Prompt) bool {
	p.buffer.DeleteBeforeCursor(1, p.renderer.col, p.renderer.row)
	return true
}

// GoRightChar Forward one character
func GoRightChar(p *Prompt) bool {
	return p.CursorRight(1)
}

// GoLeftChar Backward one character
func GoLeftChar(p *Prompt) bool {
	return p.CursorLeft(1)
}

func DeleteWordBeforeCursor(p *Prompt) bool {
	p.buffer.DeleteBeforeCursorRunes(
		istrings.RuneCountInString(p.buffer.Document().GetWordBeforeCursorWithSpace()),
		p.renderer.col,
		p.renderer.row,
	)
	return true
}

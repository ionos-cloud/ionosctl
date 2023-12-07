package prompt

import (
	istrings "github.com/elk-language/go-prompt/strings"
)

// History stores the texts that are entered.
type History struct {
	histories []string
	tmp       []string
	selected  int
}

// Add to add text in history.
func (h *History) Add(input string) {
	h.histories = append(h.histories, input)
	h.Clear()
}

// Clear to clear the history.
func (h *History) Clear() {
	h.tmp = make([]string, len(h.histories))
	copy(h.tmp, h.histories)
	h.tmp = append(h.tmp, "")
	h.selected = len(h.tmp) - 1
}

// Older saves a buffer of current line and get a buffer of previous line by up-arrow.
// The changes of line buffers are stored until new history is created.
func (h *History) Older(buf *Buffer, columns istrings.Width, rows int) (new *Buffer, changed bool) {
	if len(h.tmp) == 1 || h.selected == 0 {
		return buf, false
	}
	h.tmp[h.selected] = buf.Text()

	h.selected--
	new = NewBuffer()
	new.InsertTextMoveCursor(h.tmp[h.selected], columns, rows, false)
	return new, true
}

// Newer saves a buffer of current line and get a buffer of next line by up-arrow.
// The changes of line buffers are stored until new history is created.
func (h *History) Newer(buf *Buffer, columns istrings.Width, rows int) (new *Buffer, changed bool) {
	if h.selected >= len(h.tmp)-1 {
		return buf, false
	}
	h.tmp[h.selected] = buf.Text()

	h.selected++
	new = NewBuffer()
	new.InsertTextMoveCursor(h.tmp[h.selected], columns, rows, false)
	return new, true
}

// NewHistory returns new history object.
func NewHistory() *History {
	return &History{
		histories: []string{},
		tmp:       []string{""},
		selected:  0,
	}
}

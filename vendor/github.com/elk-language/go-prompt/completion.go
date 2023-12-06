package prompt

import (
	"strings"

	"github.com/elk-language/go-prompt/debug"
	istrings "github.com/elk-language/go-prompt/strings"
	runewidth "github.com/mattn/go-runewidth"
)

const (
	shortenSuffix = "..."
	leftPrefix    = " "
	leftSuffix    = " "
	rightPrefix   = " "
	rightSuffix   = " "
)

// Suggest represents a single suggestion
// in the auto-complete box.
type Suggest struct {
	Text        string
	Description string
}

// CompletionManager manages which suggestion is now selected.
type CompletionManager struct {
	selected       int // -1 means nothing is selected.
	tmp            []Suggest
	max            uint16
	completer      Completer
	startCharIndex istrings.RuneNumber // index of the first char of the text that should be replaced by the selected suggestion
	endCharIndex   istrings.RuneNumber // index of the last char of the text that should be replaced by the selected suggestion
	shouldUpdate   bool

	verticalScroll int
	wordSeparator  string
	showAtStart    bool
}

// GetSelectedSuggestion returns the selected item.
func (c *CompletionManager) GetSelectedSuggestion() (s Suggest, ok bool) {
	if c.selected == -1 || c.selected >= len(c.tmp) {
		return Suggest{}, false
	} else if c.selected < -1 {
		debug.Assert(false, "must not reach here")
		c.selected = -1
		return Suggest{}, false
	}

	return c.tmp[c.selected], true
}

// GetSuggestions returns the list of suggestion.
func (c *CompletionManager) GetSuggestions() []Suggest {
	return c.tmp
}

// Unselect the currently selected suggestion.
func (c *CompletionManager) Reset() {
	c.selected = -1
	c.verticalScroll = 0
	c.Update(*NewDocument())
}

// Update the suggestions.
func (c *CompletionManager) Update(in Document) {
	c.tmp, c.startCharIndex, c.endCharIndex = c.completer(in)
}

// Select the previous suggestion item.
func (c *CompletionManager) Previous() {
	if c.verticalScroll == c.selected && c.selected > 0 {
		c.verticalScroll--
	}
	c.selected--
	c.update()
}

// Next to select the next suggestion item.
func (c *CompletionManager) Next() int {
	if c.verticalScroll+int(c.max)-1 == c.selected {
		c.verticalScroll++
	}
	c.selected++
	c.update()
	return c.selected
}

// Completing returns true when the CompletionManager selects something.
func (c *CompletionManager) Completing() bool {
	return c.selected != -1
}

func (c *CompletionManager) update() {
	max := int(c.max)
	if len(c.tmp) < max {
		max = len(c.tmp)
	}

	if c.selected >= len(c.tmp) {
		c.selected = -1
		c.verticalScroll = 0
	} else if c.selected < -1 {
		c.selected = len(c.tmp) - 1
		c.verticalScroll = len(c.tmp) - max
	}
}

func deleteBreakLineCharacters(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	return s
}

func formatTexts(o []string, max istrings.Width, prefix, suffix string) (new []string, width istrings.Width) {
	l := len(o)
	n := make([]string, l)

	lenPrefix := istrings.GetWidth(prefix)
	lenSuffix := istrings.GetWidth(suffix)
	lenShorten := istrings.GetWidth(shortenSuffix)
	min := lenPrefix + lenSuffix + lenShorten
	for i := 0; i < l; i++ {
		o[i] = deleteBreakLineCharacters(o[i])

		w := istrings.GetWidth(o[i])
		if width < w {
			width = w
		}
	}

	if width == 0 {
		return n, 0
	}
	if min >= max {
		return n, 0
	}
	if lenPrefix+width+lenSuffix > max {
		width = max - lenPrefix - lenSuffix
	}

	for i := 0; i < l; i++ {
		x := istrings.GetWidth(o[i])
		if x <= width {
			spaces := strings.Repeat(" ", int(width-x))
			n[i] = prefix + o[i] + spaces + suffix
		} else if x > width {
			x := runewidth.Truncate(o[i], int(width), shortenSuffix)
			// When calling runewidth.Truncate("您好xxx您好xxx", 11, "...") returns "您好xxx..."
			// But the length of this result is 10. So we need fill right using runewidth.FillRight.
			n[i] = prefix + runewidth.FillRight(x, int(width)) + suffix
		}
	}
	return n, lenPrefix + width + lenSuffix
}

func formatSuggestions(suggests []Suggest, max istrings.Width) (new []Suggest, width istrings.Width) {
	num := len(suggests)
	new = make([]Suggest, num)

	left := make([]string, num)
	for i := 0; i < num; i++ {
		left[i] = suggests[i].Text
	}
	right := make([]string, num)
	for i := 0; i < num; i++ {
		right[i] = suggests[i].Description
	}

	left, leftWidth := formatTexts(left, max, leftPrefix, leftSuffix)
	if leftWidth == 0 {
		return []Suggest{}, 0
	}
	right, rightWidth := formatTexts(right, max-leftWidth, rightPrefix, rightSuffix)

	for i := 0; i < num; i++ {
		new[i] = Suggest{Text: left[i], Description: right[i]}
	}
	return new, istrings.Width(leftWidth + rightWidth)
}

// Constructor option for CompletionManager.
type CompletionManagerOption func(*CompletionManager)

// Set a custom completer.
func CompletionManagerWithCompleter(completer Completer) CompletionManagerOption {
	return func(c *CompletionManager) {
		c.completer = completer
	}
}

// NewCompletionManager returns an initialized CompletionManager object.
func NewCompletionManager(max uint16, opts ...CompletionManagerOption) *CompletionManager {
	c := &CompletionManager{
		selected:       -1,
		max:            max,
		completer:      NoopCompleter,
		verticalScroll: 0,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

var _ Completer = NoopCompleter

// NoopCompleter implements a Completer function
// that always returns no suggestions.
func NoopCompleter(_ Document) ([]Suggest, istrings.RuneNumber, istrings.RuneNumber) {
	return nil, 0, 0
}

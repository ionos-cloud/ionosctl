package completions

import (
	"fmt"
	"strings"
)

type Completer struct {
	fmt.Stringer

	primaryInfo   string
	secondaryInfo []string
	infoSrc       map[string]interface{}
}

func NewCompleter(infoSrc map[string]interface{}, primaryInfo string) Completer {
	if _, ok := infoSrc[primaryInfo]; !ok {
		return Completer{}
	}

	return Completer{
		infoSrc:       infoSrc,
		primaryInfo:   fmt.Sprintf("%v", infoSrc[primaryInfo]),
		secondaryInfo: nil,
	}
}

func (c Completer) String() string {
	if c.secondaryInfo == nil {
		return c.primaryInfo
	}

	return fmt.Sprintf("%s\t %s", c.primaryInfo, strings.Join(c.secondaryInfo, " "))
}

func (c Completer) AddInfo(targetInfo string, additionalFormatting ...string) Completer {
	infoRaw, ok := c.infoSrc[targetInfo]
	if !ok {
		return Completer{}
	}

	info := fmt.Sprintf("%v", infoRaw)

	for _, format := range additionalFormatting {
		info = fmt.Sprintf(format, info)
	}

	c.secondaryInfo = append(c.secondaryInfo, info)

	return c
}

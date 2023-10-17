package completions

import (
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
)

type completionInfo struct {
	primary   string
	secondary []string
}

type Completer struct {
	hasSecondaryInfo bool

	info     []completionInfo
	infoSrcs []map[string]interface{}
}

func NewCompleter(infoSrcs []map[string]interface{}, primaryInfoKey string) Completer {
	var info []completionInfo

	for _, infoSrc := range infoSrcs {
		temp, ok := infoSrc[primaryInfoKey]
		if !ok {
			return Completer{}
		}

		info = append(info, completionInfo{
			primary:   fmt.Sprintf("%v", temp),
			secondary: nil,
		})
	}

	return Completer{
		infoSrcs:         infoSrcs,
		info:             info,
		hasSecondaryInfo: false,
	}
}

func (c Completer) ToString() []string {
	if !c.hasSecondaryInfo {
		return functional.Map(c.info, func(t completionInfo) string {
			return t.primary
		})
	}

	return functional.Map(c.info, func(t completionInfo) string {
		return fmt.Sprintf("%s\t %s", t.primary, strings.Join(t.secondary, " "))
	})
}

func (c Completer) AddInfo(targetInfoKey string, additionalFormatting ...string) Completer {
	for i, infoSrc := range c.infoSrcs {
		infoRaw, ok := infoSrc[targetInfoKey]
		if !ok {
			return Completer{}
		}

		info := fmt.Sprintf("%v", infoRaw)
		for _, format := range additionalFormatting {
			info = fmt.Sprintf(format, info)
		}

		c.info[i].secondary = append(c.info[i].secondary, info)
	}

	if !c.hasSecondaryInfo {
		c.hasSecondaryInfo = true
	}

	return c
}

package routing

import (
	"regexp"
	"strings"
)

var (
	normalStrReg    = regexp.MustCompile(`^[\w\.-]+$`)
	separatedStrReg = regexp.MustCompile(`^[\w\.\-][\w\.\-\|]+[\w\.\-]$`)
	nameStrReg      = regexp.MustCompile(`^\:\w+\b`)
	surroundStrReg  = regexp.MustCompile(`^\(.+\)$`)
)

type option struct {
	name string
	str  string
	reg  *regexp.Regexp
}

func parse(frag string) []option {
	if normalStrReg.MatchString(frag) {
		return []option{option{str: frag}}
	}

	if separatedStrReg.MatchString(frag) {
		separated := strings.Split(frag, "|")
		options := make([]option, len(separated))

		for i, s := range separated {
			options[i].str = s
		}

		return options
	}

	var name string

	nameStrReg.ReplaceAllStringFunc(frag, func(n string) string {
		name = n
		return ""
	})

	if len(frag) == 0 {
		return []option{option{name: name}}
	}

	if surroundStrReg.MatchString(frag) {
		frag = frag[1 : len(frag)-1]

		if separatedStrReg.MatchString(frag) {
			separated := strings.Split(frag, "|")
			options := make([]option, len(separated))

			for i, s := range separated {
				options[i].name = name
				options[i].str = s
			}

			return options
		}

		return []option{option{name: name, reg: regexp.MustCompile(frag)}}
	}

	return nil
}

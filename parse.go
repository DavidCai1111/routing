package routing

import (
	"fmt"
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
	reg  string
}

func parse(frag string) []option {
	if frag == "" || normalStrReg.MatchString(frag) {
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

	frag = nameStrReg.ReplaceAllStringFunc(frag, func(n string) string {
		name = n[1:]
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

		return []option{option{name: name, reg: regexp.MustCompile(frag).String()}}
	}

	panic(fmt.Sprintf("routing: Invalid frag: %v", frag))
}

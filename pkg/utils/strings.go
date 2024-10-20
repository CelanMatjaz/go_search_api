package utils

import "strings"

func TrimInside(str string) string {
	for strings.Contains(str, "\n") {
		str = strings.ReplaceAll(str, "\n", " ")
	}

	for strings.Contains(str, "  ") {
		str = strings.ReplaceAll(str, "  ", " ")
	}

	for strings.Contains(str, "( ") {
		str = strings.ReplaceAll(str, "( ", "(")
	}

	for strings.Contains(str, " )") {
		str = strings.ReplaceAll(str, " )", ")")
	}

	return strings.Trim(str, " \n\t")
}

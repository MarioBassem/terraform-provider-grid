package converter

import (
	"strings"
	"unicode"
)

func ToSnakeCase(s string) string {
	ret := strings.Builder{}
	pushUnderscore := false
	if len(s) == 0 {
		return ret.String()
	}
	ret.WriteByte(byte(unicode.ToLower(rune(s[0]))))
	for i := 1; i < len(s); i++ {
		if unicode.IsUpper(rune(s[i])) {
			if pushUnderscore {
				ret.WriteByte('_')
				pushUnderscore = false
			}
		} else {
			pushUnderscore = true
		}
		ret.WriteByte(byte(unicode.ToLower(rune(s[i]))))
	}
	return ret.String()
}

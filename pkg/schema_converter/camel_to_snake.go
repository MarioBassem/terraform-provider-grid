package converter

import (
	"strings"
	"unicode"
)

func ToSnakeCase(s string) string {
	ret := strings.Builder{}
	if len(s) == 0 {
		return ret.String()
	}
	ret.WriteByte(byte(unicode.ToLower(rune(s[0]))))
	for i := 1; i < len(s); i++ {
		if (unicode.IsLower(rune(s[i-1])) && unicode.IsUpper(rune(s[i]))) ||
			(i < len(s)-1 && unicode.IsUpper(rune(s[i])) && unicode.IsLower(rune(s[i+1]))) {
			ret.WriteByte('_')
		}
		ret.WriteByte(byte(unicode.ToLower(rune(s[i]))))

	}
	return ret.String()
}

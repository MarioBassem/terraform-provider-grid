package converter

import (
	"log"
	"strings"
	"unicode"
)

func ToCamelCase(s string) string {
	log.Printf("before camel: %s", s)
	ret := strings.Builder{}
	cap := true
	for i := range s {
		if s[i] == '_' {
			cap = true
			continue
		}
		var b byte
		if cap {
			b = byte(unicode.ToUpper(rune(s[i])))
			cap = false
		} else {
			b = s[i]
		}
		ret.WriteByte(b)
	}
	log.Printf("after camel: %s", ret.String())
	return ret.String()
}

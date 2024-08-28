package textbox

import (
	"unicode"
	"unicode/utf8"
)

// This function is applicable when the cursor is at the end of the string

func UpdateInputText(text *string, r rune) bool {
	if unicode.IsControl(r) {
		// only accept backspace for now
		switch r {
		case '\b':
			if *text == "" {
				return false
			}
			_, size := utf8.DecodeLastRuneInString(*text)
			*text = (*text)[:len(*text)-size]
			return true
		default:
			return false
		}
	}
	*text = *text + string(r)
	return true
}

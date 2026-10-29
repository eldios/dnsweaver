package provider

import (
	"strconv"
	"strings"
)

// DecodeTXTStrings concatenates a TXT value written as a sequence of quoted
// character strings (RFC 1035 section 5.1), such as `"v=spf1 -all"` or
// `"part one" "part two"`. It reports false when value does not start with a
// quote or is not a well-formed sequence, in which case callers keep the
// value as given.
func DecodeTXTStrings(value string) (string, bool) {
	if !strings.HasPrefix(value, "\"") {
		return "", false
	}
	var content strings.Builder
	rest := value
	for rest != "" {
		if rest[0] != '"' {
			return "", false
		}
		end := 1
		for end < len(rest) {
			if rest[end] == '\\' {
				end += 2
				continue
			}
			if rest[end] == '"' {
				break
			}
			end++
		}
		if end >= len(rest) {
			return "", false
		}
		decoded, err := strconv.Unquote(rest[:end+1])
		if err != nil {
			return "", false
		}
		content.WriteString(decoded)
		rest = strings.TrimSpace(rest[end+1:])
	}
	return content.String(), true
}

// canonicalTXT returns the TXT value with any quoted-string presentation
// removed, so a declared, a listed and a recorded form of one value compare
// equal.
func canonicalTXT(value string) string {
	if decoded, ok := DecodeTXTStrings(value); ok {
		return decoded
	}
	return value
}

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

// TXTPayload returns the character data a declared TXT value stands for. A
// value written in quoted presentation, such as `"v=spf1 -all"` or
// `"\"hello\""`, is decoded once; any other value already is the data.
// Record.Target holds TXT data in this decoded form throughout dnsweaver.
func TXTPayload(value string) string {
	if decoded, ok := DecodeTXTStrings(value); ok {
		return decoded
	}
	return value
}

// TXTPresentation is the inverse of TXTPayload for providers whose API takes
// presentation text: data that TXTPayload would decode again, such as
// `"hello"` with its quotes, is written as quoted strings of at most 255
// bytes, and any other data is returned unchanged.
func TXTPresentation(payload string) string {
	if _, ok := DecodeTXTStrings(payload); !ok {
		return payload
	}
	var parts []string
	for len(payload) > 255 {
		parts = append(parts, strconv.Quote(payload[:255]))
		payload = payload[255:]
	}
	return strings.Join(append(parts, strconv.Quote(payload)), " ")
}

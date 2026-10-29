package provider

import (
	"strings"
	"testing"
)

func TestDecodeTXTStrings(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   string
		wantOK bool
	}{
		{"one string", `"v=spf1 -all"`, "v=spf1 -all", true},
		{"split strings", `"p=MIIB" "IjAN"`, "p=MIIBIjAN", true},
		{"escaped quote", `"say \"hi\""`, `say "hi"`, true},
		{"unquoted", "v=spf1 -all", "", false},
		{"unterminated", `"open`, "", false},
		{"text after a string", `"a" b`, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := DecodeTXTStrings(tt.input)
			if got != tt.want || ok != tt.wantOK {
				t.Errorf("DecodeTXTStrings(%q) = %q, %v; want %q, %v", tt.input, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestTXTPayloadAndPresentation(t *testing.T) {
	tests := []struct {
		name         string
		declared     string
		payload      string
		presentation string
	}{
		{"plain", "v=spf1 -all", "v=spf1 -all", "v=spf1 -all"},
		{"quoted", `"v=spf1 -all"`, "v=spf1 -all", "v=spf1 -all"},
		{"split", `"p=MIIB" "IjAN"`, "p=MIIBIjAN", "p=MIIBIjAN"},
		{"data containing quotes", `"\"hello\""`, `"hello"`, `"\"hello\""`},
		{"malformed quoting is data", `"open`, `"open`, `"open`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TXTPayload(tt.declared); got != tt.payload {
				t.Errorf("TXTPayload(%q) = %q, want %q", tt.declared, got, tt.payload)
			}
			if got := TXTPresentation(tt.payload); got != tt.presentation {
				t.Errorf("TXTPresentation(%q) = %q, want %q", tt.payload, got, tt.presentation)
			}
			if got := TXTPayload(TXTPresentation(tt.payload)); got != tt.payload {
				t.Errorf("round trip of %q = %q", tt.payload, got)
			}
		})
	}
}

func TestMemberOwnership_TXTPayloadIsExact(t *testing.T) {
	txt := func(target string) Record { return Record{Type: RecordTypeTXT, Target: target} }
	quotedData := MakeMemberOwnershipValue("x", txt(`"hello"`), nil)
	if MatchesMemberOwnership(quotedData, "x", txt("hello")) {
		t.Error("marker for the data \"hello\" matched the data hello")
	}
	if !MatchesMemberOwnership(quotedData, "x", txt(`"hello"`)) {
		t.Error("marker did not match its own data")
	}
	// Markers from before the payload flag recorded the declared spelling.
	legacy := strings.Replace(MakeMemberOwnershipValue("x", txt(`"v=spf1 -all"`), nil), ",record-txt-payload=1", "", 1)
	if !MatchesMemberOwnership(legacy, "x", txt("v=spf1 -all")) {
		t.Error("legacy marker with a quoted declaration lost its record")
	}
}

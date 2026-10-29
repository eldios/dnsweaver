package provider

import "testing"

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

func TestSameRecordMember_TXTQuoting(t *testing.T) {
	txt := func(target string) Record { return Record{Type: RecordTypeTXT, Target: target} }
	if !SameRecordMember(txt(`"v=spf1 -all"`), txt("v=spf1 -all")) {
		t.Error("quoted and unquoted forms of one value did not match")
	}
	if !SameRecordMember(txt(`"v=DKIM1; p=MIIB" "IjAN"`), txt("v=DKIM1; p=MIIBIjAN")) {
		t.Error("split and joined forms of one value did not match")
	}
	if SameRecordMember(txt(`"v=spf1 -all"`), txt("v=spf1 ~all")) {
		t.Error("different TXT values matched")
	}
	if SameRecordMember(txt(`"open`), txt("open")) {
		t.Error("malformed quoting was decoded")
	}
}

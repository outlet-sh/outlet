package emailval

import "testing"

func TestClassify(t *testing.T) {
	yes, no := true, false
	lvl := func(l Level) *Level { return &l }

	tests := []struct {
		name     string
		r        Result
		catchAll bool
		want     string
	}{
		{"bad syntax", Result{FailedAt: lvl(LevelSyntax)}, false, StatusInvalid},
		{"no MX", Result{FailedAt: lvl(LevelDomain)}, false, StatusInvalid},
		{"mailbox rejected", Result{SMTPChecked: true, SMTPDeliverable: &no}, false, StatusInvalid},
		{"rejected beats role", Result{RoleBased: true, FailedAt: lvl(LevelRole), SMTPChecked: true, SMTPDeliverable: &no}, false, StatusInvalid},
		{"disposable", Result{Disposable: true, FailedAt: lvl(LevelDisposable)}, false, StatusDisposable},
		{"role", Result{RoleBased: true, FailedAt: lvl(LevelRole)}, false, StatusRisky},
		{"catch-all", Result{SMTPChecked: true, SMTPDeliverable: &yes}, true, StatusRisky},
		{"probe error", Result{SMTPChecked: true}, false, StatusUnknown},
		{"accepted", Result{SMTPChecked: true, SMTPDeliverable: &yes}, false, StatusValid},
		{"no smtp check", Result{}, false, StatusValid},
	}
	for _, tt := range tests {
		if got := Classify(&tt.r, tt.catchAll); got != tt.want {
			t.Errorf("%s: got %s, want %s", tt.name, got, tt.want)
		}
	}
}

package entity

import "testing"

func TestUsername_Valid(t *testing.T) {
	// n := Username("dd")
	tests := []struct {
		name string
		n    UserIdentifier
		want bool
	}{
		// TODO: Add test cases.
		{"a", "a", true},
		{"", "", false},
		{"_", "_", true},
		{"asidfj;sod", "asidfj;sod", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Valid(); got != tt.want {
				t.Errorf("Username.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

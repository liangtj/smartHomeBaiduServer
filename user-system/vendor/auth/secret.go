package auth

import (
	"strings"

	"config"
)

// Secret as the type of the auth-system
type Secret string

// Verify ...
func (s Secret) Verify(t Secret) bool {
	return s == t
}

func (s Secret) size() int {
	return len(s)
}

func (s Secret) String() string {

	if config.DebugMode() {
		return string(s)
	}

	return strings.Repeat("*", s.size())
}

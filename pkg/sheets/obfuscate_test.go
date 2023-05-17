package sheets

import (
	"testing"
)

func TestObfuscationRoundTrip(t *testing.T) {
	plain := "hello world"

	actual := deobfuscate(obfuscate(plain))

	if plain != actual {
		t.Errorf("obfuscation failed: %q != %q", plain, actual)
	}
}

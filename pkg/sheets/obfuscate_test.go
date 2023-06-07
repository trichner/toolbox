package sheets

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObfuscate(t *testing.T) {
	plain := "hello world"

	s := obfuscate(plain)
	assert.Equal(t, "4ca91182c6f6d0ba7fb13e", s)
}

func TestDeobfuscate(t *testing.T) {
	o := "4ca91182c6f6d0ba7fb13e"

	s := deobfuscate(o)
	assert.Equal(t, "hello world", s)
}

func TestObfuscationRoundTrip(t *testing.T) {
	for i := 0; i < 32; i++ {
		t.Run(fmt.Sprintf("%s_%d", t.Name(), i), func(t *testing.T) {
			plain := randtext(int64(i))
			actual := deobfuscate(obfuscate(plain))
			assert.Equal(t, plain, actual)
		})
	}
}

func randtext(seed int64) string {
	r := rand.New(rand.NewSource(seed))

	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	special := "0123456789!@#$%^&*()_+-=[]{}\\;':\"|,./<>?"

	all := alphabet + strings.ToLower(alphabet) + special

	l := r.Intn(128)

	var buf strings.Builder
	for i := 0; i < l; i++ {
		buf.WriteByte(all[r.Intn(len(all))])
	}
	return buf.String()
}

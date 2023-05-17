package sheets

import (
	"encoding/hex"
	"math/rand"
)

const k = 7673551453341823650

func obfuscate(s string) string {

	plain := []byte(s)
	xor := make([]byte, len(plain))

	r := rand.New(rand.NewSource(k))
	r.Read(xor)

	for i := range plain {
		plain[i] ^= xor[i]
	}

	return hex.EncodeToString(plain)
}

func deobfuscate(s string) string {

	obfuscated, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	xor := make([]byte, len(obfuscated))

	r := rand.New(rand.NewSource(k))
	r.Read(xor)

	for i := range obfuscated {
		obfuscated[i] ^= xor[i]
	}

	return string(obfuscated)
}

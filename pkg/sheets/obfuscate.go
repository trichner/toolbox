package sheets

import (
	"encoding/hex"
	"math/rand"
)

const k = 7673551453341823650

func obfuscate(s string) string {
	plain := []byte(s)
	pad := genpad(len(plain))
	xor(plain, pad)

	return hex.EncodeToString(plain)
}

func deobfuscate(s string) string {
	obfuscated, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	pad := genpad(len(obfuscated))
	xor(obfuscated, pad)
	return string(obfuscated)
}

func xor(buf []byte, pad []byte) {
	for i := range buf {
		buf[i] ^= pad[i]
	}
}

func genpad(l int) []byte {
	pad := make([]byte, l)
	r := rand.New(rand.NewSource(k))
	r.Read(pad)
	return pad
}

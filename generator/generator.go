package generator

import (
	"math/rand"
	"time"
)

const alph = "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"

//const alph = "abc"
type Generator struct {
}

func (g Generator) ShortURL() string {
	rand.Seed(time.Now().UTC().UnixNano())
	ans := ""

	for i := 0; i < 6; i++ {
		bytes := rand.Intn(len(alph))
		ans += string(alph[bytes])
	}
	return ans
}

package key

type KeyGenerator struct {
	keyChar []byte
}

func NewKeyGenerator(keyChar []byte) *KeyGenerator {
	return &KeyGenerator{keyChar: keyChar}
}

func (g KeyGenerator) Generate(n int) string {
	if n == 0 {
		return string(g.keyChar[0])
	}

	l := len(g.keyChar)
	s := make([]byte, 20)
	i := len(s)

	for n > 0 && i >= 0 {
		i--
		j := n % l
		n = (n - j) / l
		s[i] = g.keyChar[j]
	}

	return string(s[i:])
}

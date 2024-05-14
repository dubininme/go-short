package key

type Generator struct {
	base     int
	alphabet []byte
}

func NewGenerator(alphabet []byte) (*Generator, error) {
	if len(alphabet) < 16 {
		return nil, errors.New("invalid alphabet: must be min 16 characters long")
	}

	if !isUnique(alphabet) {
		return nil, errors.New("alphabet contains repeated characters")
	}

	return &Generator{alphabet: alphabet, base: len(alphabet)}, nil
}

// Generate converts an integer to a custom hexadecimal string using a specified alphabet.
func (g *Generator) Generate(n int) string {
	if n == 0 {
		return string(g.alphabet[0])
	}

	var result strings.Builder

	// generate str
	for n > 0 {
		index := n % g.base
		result.WriteByte(g.alphabet[index])
		n /= g.base
	}

	// The result needs to be reversed as the least significant digit comes first
	return g.reverse(result.String())
}

func (g *Generator) reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// isUnique checks whether the array contains duplicate elements.
func isUnique(alphabet []byte) bool {
	seen := make(map[byte]bool)
	for _, value := range alphabet {
		if seen[value] {
			return false
		}
		seen[value] = true
	}
	return true
}

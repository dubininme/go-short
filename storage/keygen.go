package storage

type Generator interface {
	Generate(n int) string
}

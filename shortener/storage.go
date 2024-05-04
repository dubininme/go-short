package shortener

type Storage interface {
	Get(key string) string
	Put(url string) string
}

package shortener

type Storage interface {
	Get(key, url *string) error
	Put(url, key *string) error
}

package service

type Storage interface {
	Lookup(shortURL string) (string, error)
	Add(shortURL string, originalURL string) error
	Ping() error
	Close()
}

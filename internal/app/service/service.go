package service

type Storage interface {
	Get(shortURL string) (string, error)
	Insert(shortURL string, originalURL string) error
	Ping() error
	Close()
}

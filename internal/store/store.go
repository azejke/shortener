package store

type IStore interface {
	Get(id string) (string, bool)
	Insert(key string, value string)
}

type Store map[string]string

func (s *Store) Get(id string) (string, bool) {
	value, ok := (*s)[id]
	return value, ok
}

func (s *Store) Insert(key string, value string) {
	(*s)[key] = value
}

func InitStore() *Store {
	return &Store{}
}

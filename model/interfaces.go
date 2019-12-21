package model

//Storage
type Storage interface {
	Status() error
	Insert([]byte) (string, error)
	Select(string) ([]byte, error)
	SelectAll(int, int) ([]byte, error)
	Update(string, []byte) error
	Delete(string) error
	Purge(string) error
}

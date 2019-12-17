package model

//Storage
type Storage interface {
	Status() error
	Insert(Cat) (string, error)
	Select(string) (Cat, error)
	SelectAll(int, int) ([]Cat, error)
	Update(string, Cat) error
	Delete(string) error
	Purge(string) error // deletes all items from table
}

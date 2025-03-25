package models

const (
	FileCollection = "file"
	UserCollection = "user"
)

type DB interface {
	Insert(string, string, interface{}) error
	Get(string, string, interface{}) error
}

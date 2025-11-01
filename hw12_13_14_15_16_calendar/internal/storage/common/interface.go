package common

type StorageDriverInterface interface {
	Add(event Event) (Event, error)
	Update(event Event) error
	Delete(id interface{}) error
	GetByID(id interface{}) (Event, error)
	List() ([]Event, error)
	PrepareStorage(log LoggerInterface) error
}

type LoggerInterface interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

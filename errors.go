package main;


// DataNotFound is a trivial implementation of error.
type DataNotFound struct {
	message string
}

func (e *DataNotFound) Error() string {
	return e.message
}

// NotFoundException generates a new DataNotFound error
func NotFoundException(description string) error {
	return &DataNotFound{description}
}
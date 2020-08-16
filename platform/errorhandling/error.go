package errorhandling

import (
	"errors"
	"fmt"
)

// Error represents the structure to handle errors
type Error struct {
	Status  int
	Context string
	Message error `json:"error"`
}

// New returns a new instance of Error struct
func New(status int, methodName, message string) *Error {

	return &Error{
		Status:  status,
		Context: methodName,
		Message: errors.New(message),
	}

}

func (e *Error) String() string {
	return fmt.Sprintf("%s: %v. Status: %d", e.Context, e.Message, e.Status)
}

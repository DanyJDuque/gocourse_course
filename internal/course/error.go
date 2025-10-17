package course

import (
	"errors"
	"fmt"
)

var ErrNameRequiered = errors.New("name is requiered")
var ErrStartDateRequiered = errors.New("star date is requiered")
var ErrEndDateRequiered = errors.New("end date is requiered")
var ErrInvalidStartDate = errors.New("invalid start date")
var ErrInvalidEndDate = errors.New("invalid end date")

type ErrNotFound struct {
	courseID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user '%s' doesn't exist", e.courseID)
}

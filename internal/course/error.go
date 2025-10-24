package course

import (
	"errors"
	"fmt"
)

var ErrNameRequired = errors.New("name is required")
var ErrStartDateRequired = errors.New("star date is required")
var ErrEndDateRequired = errors.New("end date is required")
var ErrInvalidStartDate = errors.New("invalid start date")
var ErrInvalidEndDate = errors.New("invalid end date")
var ErrEndLesserStart = errors.New("end date mustn't be lesser than start date")

type ErrNotFound struct {
	courseID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("course '%s' doesn't exists", e.courseID)
}

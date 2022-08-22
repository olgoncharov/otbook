package errgroup

import (
	"errors"
	"strings"
)

type ErrGroup struct {
	errors    []string
	separator string
}

func NewErrorGroup(separator string) *ErrGroup {
	return &ErrGroup{
		errors:    make([]string, 0),
		separator: separator,
	}
}

func (e *ErrGroup) AddError(err error) {
	if err == nil {
		return
	}

	e.errors = append(e.errors, err.Error())
}

func (e *ErrGroup) AddErrorText(txt string) {
	if txt == "" {
		return
	}

	e.errors = append(e.errors, txt)
}

func (e *ErrGroup) Err() error {
	if len(e.errors) == 0 {
		return nil
	}

	return errors.New(strings.Join(e.errors, e.separator))
}

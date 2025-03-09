package exceptions

import "github.com/niflheimdevs/backend/internal/enums"

type Exception struct {
	Tag    enums.GeneralError
	Errors []enums.SpecificError
}

func (e *Exception) AddError(se enums.SpecificError) {
	e.Errors = append(e.Errors, se)
}
